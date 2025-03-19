package common_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/interface/controller/common"
	"github.com/kazukimurahashi12/webapp/interface/session/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetLoginIdBySession(t *testing.T) {
	t.Run("unit, happy path", func(t *testing.T) {
		// 準備
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		sessionMock := mock.NewMockSessionManager(t)
		sessionMock.On("GetSession", c).Return("id", nil)

		controller := common.NewCommonController(sessionMock)

		// 実行
		controller.GetLoginIdBySession(c)

		// 評価
		assert.Equal(t, http.StatusOK, response.Code)

		var responseJSON map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &responseJSON)
		assert.NoError(t, err)
		assert.Equal(t, "id", responseJSON["id"])
	})

	t.Run("unit, error", func(t *testing.T) {
		// 準備
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Request, _ = http.NewRequest(
			http.MethodGet,
			"/api/login-id",
			nil,
		)

		sessionMock := mock.NewMockSessionManager(t)
		sessionMock.On("GetSession", c).Return("", errors.New("session error"))

		controller := common.NewCommonController(sessionMock)

		// 実行
		controller.GetLoginIdBySession(c)

		// 評価
		assert.Equal(t, http.StatusUnauthorized, response.Code)

		var responseJSON map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &responseJSON)
		assert.NoError(t, err)
		assert.Equal(t, "session error", responseJSON["error"])
	})
}
