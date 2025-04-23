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

func TestLoginController_GetLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		// 準備
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
		ctx.Request = req

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockAuthUseCase := authMocks.NewMockUseCase(ctrl)

		// セッションモック
		mockSession.EXPECT().
			GetSession(gomock.Any()).
			Return("user123", nil)

		// ユーザーモック
		mockAuthUseCase.EXPECT().
			GetUserByID("user123").
			Return(&user.User{UserID: "user123"}, nil)

		logger := zaptest.NewLogger(t)
		controller := NewLoginController(mockAuthUseCase, mockSession, logger)

		// 実行
		controller.GetLogin(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	})

	t.Run("SessionError", func(t *testing.T) {
		// 準備
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
		ctx.Request = req

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockAuthUseCase := authMocks.NewMockUseCase(ctrl)

		mockSession.EXPECT().
			GetSession(gomock.Any()).
			Return("", errors.New("session error"))

		logger := zaptest.NewLogger(t)
		controller := NewLoginController(mockAuthUseCase, mockSession, logger)

		// 実行
		controller.GetLogin(ctx)

		// 検証
		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
	})

	t.Run("UserFetchError", func(t *testing.T) {
		// 準備
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
		ctx.Request = req

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockAuthUseCase := authMocks.NewMockUseCase(ctrl)

		mockSession.EXPECT().
			GetSession(gomock.Any()).
			Return("user123", nil)

		mockAuthUseCase.EXPECT().
			GetUserByID("user123").
			Return(nil, errors.New("user not found"))

		logger := zaptest.NewLogger(t)
		controller := NewLoginController(mockAuthUseCase, mockSession, logger)

		// 実行
		controller.GetLogin(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})
}

func TestLoginController_PostLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		// 準備
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"UserID":"testuser","Password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockAuthUseCase := authMocks.NewMockUseCase(ctrl)

		// 認証モック
		mockAuthUseCase.EXPECT().
			Authenticate("testuser", "password123").
			Return(&user.User{UserID: "user123"}, nil)

		// セッション作成モック
		mockSession.EXPECT().
			CreateSession("user123").
			Return(nil)

		logger := zaptest.NewLogger(t)
		controller := NewLoginController(mockAuthUseCase, mockSession, logger)
		// 実行
		controller.PostLogin(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ログインに成功しました", response["message"])
			assert.Equal(t, "LOGIN_SUCCESS", response["code"])
		}
	})

	t.Run("AuthenticationFailed", func(t *testing.T) {
		// 準備
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"UserID":"testuser","Password":"wrongpassword"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockAuthUseCase := authMocks.NewMockUseCase(ctrl)

		mockAuthUseCase.EXPECT().
			Authenticate("testuser", "wrongpassword").
			Return(nil, errors.New("authentication failed"))

		logger := zaptest.NewLogger(t)
		controller := NewLoginController(mockAuthUseCase, mockSession, logger)
		// 実行
		controller.PostLogin(ctx)

		// 検証
		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status())
		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ユーザーIDまたはパスワードが正しくありません", response["error"])
			assert.Equal(t, "AUTHENTICATION_FAILED", response["code"])
		}
	})

	t.Run("SessionCreationFailed", func(t *testing.T) {
		// 準備
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"UserID":"testuser","Password":"password123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockAuthUseCase := authMocks.NewMockUseCase(ctrl)

		mockAuthUseCase.EXPECT().
			Authenticate("testuser", "password123").
			Return(&user.User{UserID: "user123"}, nil)

		mockSession.EXPECT().
			CreateSession("user123").
			Return(errors.New("session creation failed"))

		logger := zaptest.NewLogger(t)
		controller := NewLoginController(mockAuthUseCase, mockSession, logger)
		// 実行
		controller.PostLogin(ctx)
		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "セッションの作成に失敗しました", response["error"])
			assert.Equal(t, "SESSION_CREATION_FAILED", response["code"])
		}
	})
}
