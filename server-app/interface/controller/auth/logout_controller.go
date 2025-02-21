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

func (l *LogoutController) DecideLogout(c *gin.Context) {
	var logoutUser domain.FormUser
	if err := c.ShouldBindJSON(&logoutUser); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := l.authUseCase.Authenticate(logoutUser.UserID, logoutUser.Password)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := l.sessionManager.DeleteSession(c); err != nil {
		log.Printf("Failed to delete session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Successfully logged out:", user.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
