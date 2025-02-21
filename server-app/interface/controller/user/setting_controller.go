package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/user"
)

type SettingController struct {
	userUseCase    user.UseCase
	sessionManager session.SessionManager
}

func NewSettingController(userUseCase user.UseCase, sessionManager session.SessionManager) *SettingController {
	return &SettingController{
		userUseCase:    userUseCase,
		sessionManager: sessionManager,
	}
}

func (s *SettingController) UpdateID(c *gin.Context) {
	var userUpdate domain.UserIdChange
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := s.userUseCase.UpdateUserID(userUpdate.ChangeID, userUpdate.NowID)
	if err != nil {
		log.Printf("Failed to update user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.sessionManager.UpdateSession(c, userUpdate.ChangeID, userUpdate.NowID); err != nil {
		log.Printf("Failed to update session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

func (s *SettingController) UpdatePassword(c *gin.Context) {
	var passwordUpdate domain.UserPwChange
	if err := c.ShouldBindJSON(&passwordUpdate); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := s.userUseCase.UpdateUserPassword(
		passwordUpdate.UserID,
		passwordUpdate.NowPassword,
		passwordUpdate.ChangePassword,
	)
	if err != nil {
		log.Printf("Failed to update password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.sessionManager.DeleteSession(c); err != nil {
		log.Printf("Failed to delete session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.sessionManager.CreateSession(updatedUser.ID); err != nil {
		log.Printf("Failed to create new session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}
