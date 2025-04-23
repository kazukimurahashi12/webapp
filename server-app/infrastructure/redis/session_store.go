package redis

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/kazukimurahashi12/webapp/interface/session"
)

//#######################################
// Redis接続処理
//#######################################

var _ session.SessionManager = &RedisSessionStore{}

type RedisSessionStore struct {
	conn *redis.Client
}

// RedisClientインスタン作成
func NewRedisSessionStore() *RedisSessionStore {
	// Zapロガー初期化
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// 環境変数設定
	// プロジェクトルートディレクトリを取得
	rootDir := os.Getenv("PROJECT_ROOT")
	if rootDir == "" {
		logger.Error("PROJECT_ROOT environment variable is not set")
		return &RedisSessionStore{conn: nil}
	}

	// 環境変数ファイルのパス
	envPath := filepath.Join(rootDir, "build", "db", "data", ".env")

	// 環境変数ファイルを読み込み
	if err := godotenv.Load(envPath); err != nil {
		logger.Error("Failed to load .env file",
			zap.String("path", envPath),
			zap.Error(err))
		return &RedisSessionStore{conn: nil}
	}

	var dbHost string
	if os.Getenv("DOCKER_ENV") == "true" {
		// Dockerコンテナ内での接続先を指定
		dbHost = os.Getenv("REDIS_DOCKER_HOST")
	} else {
		// ローカル環境での接続先を指定
		dbHost = os.Getenv("REDIS_LOCAL_HOST")
	}
	// Redisデータベース接続のためRedisクライアント作成
	conn := redis.NewClient(&redis.Options{
		Addr:     dbHost,
		Password: "",
		DB:       0,
	})
	return &RedisSessionStore{conn: conn}
}

// セッションを作成
func (s *RedisSessionStore) CreateSession(userID string) error {
	slice := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, slice); err != nil {
		log.Println("ランダムな文字作成時にエラーが発生しました。", err.Error())
		return err
	}

	redisKey := base64.URLEncoding.EncodeToString(slice)
	if err := s.conn.Set(context.Background(), redisKey, userID, 0).Err(); err != nil {
		log.Println("Session登録時にエラーが発生:", err.Error())
		return err
	}
	return nil
}

// セッションを取得
func (s *RedisSessionStore) GetSession(c *gin.Context) (string, error) {
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	redisKey, err := c.Cookie(cookieKey)
	if err != nil {
		log.Printf("Session cookie not found. redisKey: %s, cookieKey: %s, err: %v", redisKey, cookieKey, err)
		return "", err
	}

	redisValue, err := s.conn.Get(context.Background(), redisKey).Result()
	if err != nil {
		log.Printf("Failed to get session data from Redis. redisKey: %s, redisValue: %s, err: %v", redisKey, redisValue, err)
		return "", err
	}
	return redisValue, nil
}

// セッションを削除
func (s *RedisSessionStore) DeleteSession(c *gin.Context) error {
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	redisKey, err := c.Cookie(cookieKey)
	if err != nil {
		log.Println("セッションのクッキーが見つかりませんでした。,err:" + err.Error())
		return err
	}

	if err := s.conn.Del(context.Background(), redisKey).Err(); err != nil {
		log.Printf("Failed to delete session from Redis. redisKey: %s, err: %v", redisKey, err)
		return err
	}

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
	http.SetCookie(c.Writer, cookie)
	return nil
}

// セッションを更新
func (s *RedisSessionStore) UpdateSession(c *gin.Context, newID string) error {
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	redisKey, err := c.Cookie(cookieKey)
	if err != nil {
		log.Println("セッションのクッキーが見つかりませんでした。,err:" + err.Error())
		return err
	}

	// 古いセッションを削除
	if err := s.conn.Del(context.Background(), redisKey).Err(); err != nil {
		log.Printf("Failed to delete session from Redis. redisKey: %s, err: %v", redisKey, err)
		return err
	}

	// 新しいセッションを作成
	slice := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, slice); err != nil {
		log.Println("ランダムな文字作成時にエラーが発生しました。", err.Error())
		return err
	}

	newRedisKey := base64.URLEncoding.EncodeToString(slice)
	if err := s.conn.Set(context.Background(), newRedisKey, newID, 0).Err(); err != nil {
		log.Println("Session登録時にエラーが発生:", err.Error())
		return err
	}

	// 新しいクッキーを設定
	cookie := &http.Cookie{
		Name:     cookieKey,
		Value:    newRedisKey,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   0,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(c.Writer, cookie)
	return nil
}
