package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/infrastructure/redis"
	redis_mock "github.com/kazukimurahashi12/webapp/infrastructure/redis/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetLoginIdBySession(t *testing.T) {
	t.Run("integration, happy path", func(t *testing.T) {
		//GetLoginIdBySessionインテグレーションテスト、正常系（redisアクセス）
		// Arrange ---
		res1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(res1)
		c1.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		// セッション情報を設定
		os.Setenv("LOGIN_USER_ID_KEY", "loginUserIdKey") // テスト用の環境変数を設定
		redis := redis.NewRedisSessionStore()
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

		// Act ---
		common.GetLoginIdBySession(c2, redis)

		// Assert ---
		assert.Equal(t, http.StatusOK, res2.Code) // ステータスコードが期待通りか確認

		// レスポンスボディをJSONとしてパース
		var responseJSON map[string]interface{}
		err = json.Unmarshal(res2.Body.Bytes(), &responseJSON)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		// レスポンスボディの"id"フィールドを確認
		expectedID := "root" // 期待されるIDを設定
		if id, found := responseJSON["id"]; found {
			assert.Equal(t, expectedID, id)
		} else {
			t.Fatalf("Response does not contain 'id' field")
		}
	})
	t.Run("unit, happy path", func(t *testing.T) {
		//GetLoginIdBySession単体テスト、正常系・異常系（モック利用）
		// Arrange ---
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)
		redis := redis_mock.NewSessionStore(t)
		redis.On("GetSession", c, "loginUserIdKey").Return("id", nil)
		// Act ---
		common.GetLoginIdBySession(c, redis)
		// Assert ---
		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("unit, error", func(t *testing.T) {
		// Arrange ---
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)
		redis := redis_mock.NewSessionStore(t)
		redis.On("GetSession", c, "loginUserIdKey").Return("", errors.New("error"))
		// Act ---
		common.GetLoginIdBySession(c, redis)
		// Assert ---
		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})
}
