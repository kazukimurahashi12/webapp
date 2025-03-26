package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/auth"
	"go.uber.org/zap"
)

//#######################################
//ログアウトコントローラー
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
	var logoutUser domain.FormUser
	if err := c.ShouldBindJSON(&logoutUser); err != nil {
		l.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエスト形式が不正です",
			"code":  "INVALID_REQUEST_FORMAT",
		})
		return
	}

	// UseCaseユーザー認証
	user, err := l.authUseCase.Authenticate(logoutUser.UserId, logoutUser.Password)
	if err != nil {
		l.logger.Error("Authentication failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ユーザーIDまたはパスワードが正しくありません",
			"code":  "AUTHENTICATION_FAILED",
		})
		return
	}

	// セッション削除
	if err := l.sessionManager.DeleteSession(c); err != nil {
		l.logger.Error("Failed to delete session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "セッションの削除に失敗しました",
			"code":  "SESSION_DELETION_FAILED",
		})
		return
	}

	// ログアウト成功
	l.logger.Info("Successfully logged out", zap.Uint("userId", user.Id))
	c.JSON(http.StatusOK, gin.H{
		"message": "ログアウトに成功しました",
		"code":    "LOGOUT_SUCCESS",
	})
}
