package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/auth"
)

type LogoutController struct {
	authUseCase    auth.UseCase
	sessionManager session.SessionManager
}

func NewLogoutController(authUseCase auth.UseCase, sessionManager session.SessionManager) *LogoutController {
	return &LogoutController{
		authUseCase:    authUseCase,
		sessionManager: sessionManager,
	}
}

// ログアウト処理
func (l *LogoutController) DecideLogout(c *gin.Context) {
	var logoutUser domain.FormUser
	if err := c.ShouldBindJSON(&logoutUser); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// UseCaseユーザー認証
	user, err := l.authUseCase.Authenticate(logoutUser.UserId, logoutUser.Password)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// セッション削除
	if err := l.sessionManager.DeleteSession(c); err != nil {
		log.Printf("Failed to delete session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ログアウト成功
	log.Println("Successfully logged out:", user.Id)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
