package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	sessionMocks "github.com/kazukimurahashi12/webapp/interface/session/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestCommonController_GetLoginIdBySession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/common/login-id", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		logger := zaptest.NewLogger(t)
		controller := NewCommonController(mockSession, logger)

		// 実行
		controller.GetLoginIdBySession(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())

		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ログインIDを取得しました", response["message"])
			assert.Equal(t, "LOGIN_ID_FETCHED", response["code"])
			assert.Equal(t, "user123", response["loginID"])
		}
	})

	t.Run("UserIDNotFound", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/common/login-id", nil)
		ctx.Request = req
		// userIDを設定しない

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		logger := zaptest.NewLogger(t)
		controller := NewCommonController(mockSession, logger)

		// 実行
		controller.GetLoginIdBySession(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})
}
