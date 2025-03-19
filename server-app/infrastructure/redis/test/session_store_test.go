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
	t.Run("正常系", func(t *testing.T) {
		// 準備
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		// 実行
		store := redis.NewRedisSessionStore()
		err := store.CreateSession("test-user")

		// 検証
		assert.NoError(t, err)
	})
}

func TestGetSession(t *testing.T) {
	setUp := func(store *redis.RedisSessionStore) *gin.Context {
		res1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(res1)
		c1.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		// セッションを設定
		os.Setenv("LOGIN_USER_ID_KEY", "loginUserIdKey")
		err := store.CreateSession("root")
		assert.NoError(t, err)

		// クッキーを取得
		cookie, err := c1.Cookie("loginUserIdKey")
		assert.NoError(t, err)

		// クッキーを設定した新しいリクエストを作成
		res2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(res2)
		c2.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)
		c2.Request.AddCookie(&http.Cookie{
			Name:  "loginUserIdKey",
			Value: cookie,
		})
		return c2
	}

	t.Run("正常系", func(t *testing.T) {
		// 準備
		store := redis.NewRedisSessionStore()
		c := setUp(store)

		// 実行
		userID, err := store.GetSession(c)

		// 検証
		assert.NoError(t, err)
		assert.Equal(t, "root", userID)
	})

	t.Run("無効なクッキーキー", func(t *testing.T) {
		// 準備
		store := redis.NewRedisSessionStore()
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		// 実行
		_, err := store.GetSession(c)

		// 検証
		assert.Error(t, err)
	})

	t.Run("セッションキーが見つからない場合", func(t *testing.T) {
		// 準備
		store := redis.NewRedisSessionStore()
		c := setUp(store)
		err := store.DeleteSession(c)
		assert.NoError(t, err)

		// 実行
		_, err = store.GetSession(c)

		// 検証
		assert.Error(t, err)
	})
}

func TestDeleteSession(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		// 準備
		store := redis.NewRedisSessionStore()
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)
		err := store.CreateSession("test-user")
		assert.NoError(t, err)

		// 実行
		err = store.DeleteSession(c)

		// 検証
		assert.NoError(t, err)
	})
}
