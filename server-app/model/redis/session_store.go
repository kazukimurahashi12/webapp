package redis

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// mockを更新する場合は以下を実行
// go run github.com/vektra/mockery/v2@v2.36.0
type SessionStore interface {
	NewSession(c *gin.Context, cookieKey, redisValue string) error
	GetSession(c *gin.Context, cookieKey string) (string, error)
	DeleteSession(c *gin.Context, id string) error
	UpdateSession(c *gin.Context, ChangeId string, oldId string) error
}

type RedisSessionStore struct {
	conn *redis.Client
}

var _ SessionStore = &RedisSessionStore{}

func NewRedisSessionStore() *RedisSessionStore {
	//環境変数設定
	//main.goからの相対パス指定
	envErr := godotenv.Load("./build/db/data/.env")
	if envErr != nil {
		log.Println("Error loading .env file", envErr)
	}
	var dbHost string
	if os.Getenv("DOCKER_ENV") == "true" {
		// Dockerコンテナ内での接続先を指定
		dbHost = os.Getenv("REDIS_DOCKER_HOST")
	} else {
		// ローカル環境での接続先を指定
		dbHost = os.Getenv("REDIS_LOCAL_HOST")
	}
	//Redisデータベース接続のためRedisクライアント作成
	conn := redis.NewClient(&redis.Options{
		Addr:     dbHost,
		Password: "",
		DB:       0,
	})
	return &RedisSessionStore{conn: conn}
}

func (s *RedisSessionStore) NewSession(c *gin.Context, cookieKey, redisValue string) error {
	slice := make([]byte, 64)
	//ランダムなバイト列を生成
	if _, err := io.ReadFull(rand.Reader, slice); err != nil {
		log.Println("ランダムな文字作成時にエラーが発生しました。", err.Error())
		return err
	}

	//バイト配列を base64 エンコードして文字列に変換
	newRedisKey := base64.URLEncoding.EncodeToString(slice)

	//Redisにセッションを登録
	if err := s.conn.Set(c, newRedisKey, redisValue, 0).Err(); err != nil {
		log.Println("Session登録時にエラーが発生:", err.Error())
		return err
	}

	// SameSite属性をNoneにするために、Secure属性（HTTPS）を設定
	var secure bool = false
	if c.Request.URL.Scheme == "https" {
		secure = true
	}
	//HTTPレスポンスヘッダーにCookieを設定
	cookie := &http.Cookie{
		Name:     cookieKey,
		Value:    newRedisKey,
		Path:     "/",
		Domain:   "localhost", // 本番環境では正しいドメインを設定
		MaxAge:   0,
		HttpOnly: false,
		Secure:   secure, // 本番環境ではHTTPSでない場合はfalseにする
		SameSite: http.SameSiteNoneMode,
	}
	// クッキーをレスポンスに設定
	http.SetCookie(c.Writer, cookie)
	log.Println("クッキーをレスポンスに設定に成功。", cookie)
	return nil
}

func (s *RedisSessionStore) GetSession(c *gin.Context, cookieKey string) (string, error) {
	//クライアントのリクエストに含まれるセッションのクッキー値を取得
	redisKey, err := c.Cookie(cookieKey)
	if err != nil {
		log.Printf("セッションのクッキーが見つかりませんでした。redisKey: %s, cookieKey: %s, err: %v", redisKey, cookieKey, err)
		return "", err
	}

	//取得したセッションのクッキー値を使用して、Redisから対応するセッションデータを取得
	redisValue, err := s.conn.Get(c, redisKey).Result()
	if err != nil {
		log.Printf("Redisから対応するセッションデータを取得に失敗しました。redisKey: %s, redisValue: %s, err: %v", redisKey, redisValue, err)
		return "", err
	}
	switch {
	case err == redis.Nil:
		log.Println("SessionKeyが登録されていません。", err.Error())
		return "", err
	case err != nil:
		log.Println("Session取得時にエラー発生:", err.Error())
		return "", err
	}

	log.Printf("redisからセッション情報のIDを取得に成功。 ID: %s", redisValue)
	return redisValue, nil
}

func (s *RedisSessionStore) DeleteSession(c *gin.Context, id string) error {
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	redisKey, err := c.Cookie(cookieKey)
	if err != nil {
		log.Println("セッションのクッキーが見つかりませんでした。,err:" + err.Error())
		return err
	}

	//取得したセッションのクッキー値を使用して、Redisから対応するセッションデータを取得
	redisValue, err := s.conn.Get(c, redisKey).Result()
	if err != nil {
		log.Printf("Redisから対応するセッションデータを取得に失敗しました。redisKey: %s, redisValue: %s, err: %v", redisKey, redisValue, err)
		return err
	}

	//Redisからセッションを削除
	if redisValue == id {
		cmd := s.conn.Del(c, redisKey)
		if cmd.Val() == 0 {
			err := errors.New("error in redis of updateSession")
			log.Printf("Redisからセッションを削除できませんでした。cmd.Val(): %v", cmd.Val())
			return err
		} else {
			log.Println("Redisからセッションを削除しました。:", cmd.String())
		}
	} else {
		err := errors.New("error in redis of updateSession")
		log.Printf("セッションidが一致しませんでした。id: %s, redisValue: %s, err: %v", id, redisValue, err)
		return err
	}

	//クライアントのブラウザに保存されているセッションのクッキーを削除
	cookie := &http.Cookie{
		Name:     cookieKey,
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	// クッキーをレスポンスに設定
	http.SetCookie(c.Writer, cookie)
	log.Printf("クッキーを消去に成功。")
	return err
}

func (s *RedisSessionStore) UpdateSession(c *gin.Context, ChangeId string, oldId string) error {
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	//リクエストからCookie取得
	redisKey, err := c.Cookie(cookieKey)
	if err != nil {
		log.Println("セッションからIDが取得できませんでした。,err:" + err.Error())
		return err
	}

	//取得したセッションのクッキー値を使用して、Redisから対応するセッションデータを取得
	redisValue, err := s.conn.Get(c, redisKey).Result()
	if err != nil {
		log.Printf("Redisから対応するセッションデータを取得に失敗しました。redisKey: %s, redisValue: %s, err: %v", redisKey, redisValue, err)
		return err
	}

	//Redisからセッションを消去
	if redisValue == oldId {
		cmd := s.conn.Del(c, redisKey)
		if cmd.Val() == 0 {
			err := errors.New("error in redis of updateSession")
			log.Printf("Redisからセッションを削除できませんでした。cmd.Val(): %v", cmd.Val())
			return err
		} else {
			log.Println("Redisからセッションを削除しました。:", cmd.String())
		}
	} else {
		err := errors.New("error in redis of updateSession")
		log.Printf("セッションidが一致しませんでした。oldId: %s, redisValue: %s, err: %v", oldId, redisValue, err)
		return err
	}

	slice := make([]byte, 64)
	//ランダムなバイト列を生成
	if _, err := io.ReadFull(rand.Reader, slice); err != nil {
		log.Println("ランダムな文字作成時にエラーが発生しました。：" + err.Error())
		return err
	}

	//バイト配列を base64 エンコードして文字列に変換
	newRedisKey := base64.URLEncoding.EncodeToString(slice)

	//Redisにセッションを登録
	if err := s.conn.Set(c, newRedisKey, ChangeId, 0).Err(); err != nil {
		log.Println("Session登録時にエラーが発生:" + err.Error())
		return err
	}
	log.Printf("Redisでセッション登録に成功しました。cookieKey: %s, newRedisKey: %s", cookieKey, newRedisKey)

	// SameSite属性をNoneにするために、Secure属性（HTTPS）を設定
	var secure bool = false
	if c.Request.URL.Scheme == "https" {
		secure = true
	}
	//HTTPレスポンスヘッダーにCookieを設定
	cookie := &http.Cookie{
		Name:     cookieKey,
		Value:    newRedisKey,
		Path:     "/",
		Domain:   "localhost", // 本番環境では正しいドメインを設定
		MaxAge:   0,
		HttpOnly: false,
		Secure:   secure, // 本番環境ではHTTPSでない場合はfalseにする
		SameSite: http.SameSiteNoneMode,
	}
	// クッキーをレスポンスに設定
	http.SetCookie(c.Writer, cookie)
	log.Printf("クッキーをレスポンスに設定に成功。")
	return err
}
