package redis_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/infrastructure/redis"
	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		// 準備
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		// 実行
		cookieKey := "loginUserIdKey"
		redisValue := "root"
		redis := redis.NewRedisSessionStore()
		err := redis.CreateSession(c, cookieKey, redisValue)

		// 検証
		assert.NoError(t, err)
	})
}

func TestGetSession(t *testing.T) {
	setUp := func(redis redis.SessionStore) *gin.Context {
		res1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(res1)
		c1.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		// セッション情報を設定
		os.Setenv("LOGIN_USER_ID_KEY", "loginUserIdKey") // テスト用の環境変数を設定
		err := redis.NewSession(c1, "loginUserIdKey", "root")
		assert.NoError(t, err)

		// クッキーが設定されているか確認
		firstResponseCookies := res1.Result().Cookies()
		assert.NotEmpty(t, firstResponseCookies)

		// 同一リクエスト内でクッキーの値を読むことはできないため、新しいリクエストを作成
		res2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(res2)
		c2.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		// 最初のリクエストで設定されたクッキーを2回目のリクエストに設定
		for _, cookie := range firstResponseCookies {
			c2.Request.AddCookie(cookie)
		}
		return c2
	}
	t.Run("normal case", func(t *testing.T) {
		// 準備
		redis := redis.NewRedisSessionStore()
		c := setUp(redis)

		// 実行
		redisValue, err := redis.GetSession(c, "loginUserIdKey")

		// 検証
		assert.NoError(t, err)
		assert.Equal(t, "root", redisValue)
	})
	t.Run("invalid cookie key", func(t *testing.T) {
		// 準備
		redis := redis.NewRedisSessionStore()
		c := setUp(redis)

		//実行
		redisValue, err := redis.GetSession(c, "invalid")

		// 検証
		assert.EqualError(t, err, "http: named cookie not present")
		assert.Empty(t, redisValue)
	})
	t.Run("session key not found", func(t *testing.T) {
		//準備
		redis := redis.NewRedisSessionStore()
		c := setUp(redis)
		redis.DeleteSession(c, "root")

		// 実行
		redisValue, err := redis.GetSession(c, "loginUserIdKey")

		// 検証
		assert.EqualError(t, err, "redis: nil")
		assert.Empty(t, redisValue)
	})
}
