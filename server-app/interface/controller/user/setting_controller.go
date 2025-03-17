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

// 会員情報編集(id)
func (s *SettingController) UpdateID(c *gin.Context) {
	var userUpdate domain.UserIdChange
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// UpdateUserID処理UseCase
	updatedUser, err := s.userUseCase.UpdateUserID(userUpdate.ChangeId, userUpdate.NowId)
	if err != nil {
		log.Printf("Failed to update user ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//redisでセッション破棄、新IDでセッション作成
	if err := s.sessionManager.UpdateSession(c, userUpdate.ChangeId, userUpdate.NowId); err != nil {
		log.Printf("Failed to update session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Change UserId :userUpdate.ChangeId %+v", userUpdate.ChangeId)
	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

// 会員情報編集(password)
func (s *SettingController) UpdatePassword(c *gin.Context) {
	var passwordUpdate domain.UserPwChange
	if err := c.ShouldBindJSON(&passwordUpdate); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO 認証チェック

	// UpdateUserPassword処理UseCase
	updatedUser, err := s.userUseCase.UpdateUserPassword(
		passwordUpdate.UserId,
		passwordUpdate.NowPassword,
		passwordUpdate.ChangePassword,
	)
	if err != nil {
		log.Printf("Failed to update password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Redisよりログイン情報セッションを一度消去
	if err := s.sessionManager.DeleteSession(c); err != nil {
		log.Printf("Failed to delete session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//RedisよりセッションとCookieにUserIdを新しく登録
	if err := s.sessionManager.CreateSession(updatedUser.UserId); err != nil {
		log.Printf("Failed to create new session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//会員情報のPW変更に成功
	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}
