package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kazukimurahashi12/webapp/domain/user"
	sessionMocks "github.com/kazukimurahashi12/webapp/interface/session/mocks"
	authMocks "github.com/kazukimurahashi12/webapp/usecase/auth/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestLogoutController_DecideLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		// 準備
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"UserID":"testuser","Password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockAuthUseCase := authMocks.NewMockUseCase(ctrl)

		// 認証モック
		mockAuthUseCase.EXPECT().
			Authenticate("testuser", "password123").
			Return(&user.User{UserID: "user123"}, nil)

		// セッション削除モック
		mockSession.EXPECT().
			DeleteSession(gomock.Any()).
			Return(nil)

		logger := zaptest.NewLogger(t)
		controller := NewLogoutController(mockAuthUseCase, mockSession, logger)

		// 実行
		controller.DecideLogout(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())

		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ログアウトに成功しました", response["message"])
			assert.Equal(t, "LOGOUT_SUCCESS", response["code"])
		}
	})

	t.Run("AuthenticationFailed", func(t *testing.T) {
		// 準備
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"UserID":"testuser","Password":"wrongpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockAuthUseCase := authMocks.NewMockUseCase(ctrl)

		// 認証失敗モック
		mockAuthUseCase.EXPECT().
			Authenticate("testuser", "wrongpassword").
			Return(nil, errors.New("authentication failed"))

		logger := zaptest.NewLogger(t)
		controller := NewLogoutController(mockAuthUseCase, mockSession, logger)

		// 実行
		controller.DecideLogout(ctx)

		// 検証
		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())

		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ユーザーIDまたはパスワードが正しくありません", response["error"])
			assert.Equal(t, "AUTHENTICATION_FAILED", response["code"])
		}
	})

	t.Run("SessionDeletionFailed", func(t *testing.T) {
		// 準備
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"UserID":"testuser","Password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockAuthUseCase := authMocks.NewMockUseCase(ctrl)

		// 認証モック
		mockAuthUseCase.EXPECT().
			Authenticate("testuser", "password123").
			Return(&user.User{UserID: "user123"}, nil)

		// セッション削除失敗モック
		mockSession.EXPECT().
			DeleteSession(gomock.Any()).
			Return(errors.New("session deletion failed"))

		logger := zaptest.NewLogger(t)
		controller := NewLogoutController(mockAuthUseCase, mockSession, logger)

		// 実行
		controller.DecideLogout(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())

		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "セッションの削除に失敗しました", response["error"])
			assert.Equal(t, "SESSION_DELETION_FAILED", response["code"])
		}
	})
}
