package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/infrastructure/web/middleware"
	"github.com/kazukimurahashi12/webapp/interface/dto"
	"github.com/kazukimurahashi12/webapp/interface/mapper"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/auth"
	"go.uber.org/zap"
)

//#######################################
// ログアウトコントローラー
//#######################################

type LogoutController struct {
	authUseCase    auth.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewLogoutController(authUseCase auth.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *LogoutController {
	return &LogoutController{
		authUseCase:    authUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// ログアウト処理
func (l *LogoutController) DecideLogout(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	var logoutUser dto.FormUser
	if err := c.ShouldBindJSON(&logoutUser); err != nil {
		l.logger.Error("Failed to bind JSON",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "リクエスト形式が不正です",
			"code":       "INVALID_REQUEST_FORMAT",
			"request_id": requestID,
		})
		return
	}

	// UseCaseユーザー認証
	user, err := l.authUseCase.Authenticate(logoutUser.UserID, logoutUser.Password)
	if err != nil {
		l.logger.Error("Authentication failed",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "ユーザーIDまたはパスワードが正しくありません",
			"code":       "AUTHENTICATION_FAILED",
			"request_id": requestID,
		})
		return
	}

	// セッション削除
	if err := l.sessionManager.DeleteSession(c); err != nil {
		l.logger.Error("Failed to delete session",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "セッションの削除に失敗しました",
			"code":       "SESSION_DELETION_FAILED",
			"request_id": requestID,
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToUserIDResponse(user)
	// ログアウト成功
	l.logger.Info("Successfully logged out",
		zap.String("requestID", requestID),
		zap.Any("userId", response))
	c.JSON(http.StatusOK, gin.H{
		"message":    "ログアウトに成功しました",
		"code":       "LOGOUT_SUCCESS",
		"request_id": requestID,
		"user":       response,
	})
}
