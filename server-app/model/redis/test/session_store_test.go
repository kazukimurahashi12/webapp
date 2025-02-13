package redis_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/kazukimurahashi12/webapp/model/redis"
)

func TestNewSession(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		// Arrange ---
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		// Act ---
		cookieKey := "loginUserIdKey"
		redisValue := "root"
		redis := redis.NewRedisSessionStore()
		err := redis.NewSession(c, cookieKey, redisValue)

		// Assert ---
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
		// Arrange ---
		redis := redis.NewRedisSessionStore()
		c := setUp(redis)

		// Act ---
		redisValue, err := redis.GetSession(c, "loginUserIdKey")

		// Assert ---
		assert.NoError(t, err)
		assert.Equal(t, "root", redisValue)
	})
	t.Run("invalid cookie key", func(t *testing.T) {
		// Arrange ---
		redis := redis.NewRedisSessionStore()
		c := setUp(redis)

		// Act ---
		redisValue, err := redis.GetSession(c, "invalid")

		// Assert ---
		assert.EqualError(t, err, "http: named cookie not present")
		assert.Empty(t, redisValue)
	})
	t.Run("session key not found", func(t *testing.T) {
		// Arrange ---
		redis := redis.NewRedisSessionStore()
		c := setUp(redis)
		redis.DeleteSession(c, "root")

		// Act ---
		redisValue, err := redis.GetSession(c, "loginUserIdKey")

		// Assert ---
		assert.EqualError(t, err, "redis: nil")
		assert.Empty(t, redisValue)
	})
}
