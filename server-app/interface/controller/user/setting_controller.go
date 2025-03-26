package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/user"
	"go.uber.org/zap"
)

type SettingController struct {
	userUseCase    user.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewSettingController(userUseCase user.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *SettingController {
	return &SettingController{
		userUseCase:    userUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// 会員情報編集(id)
func (s *SettingController) UpdateID(c *gin.Context) {
	var userUpdate domain.UserIdChange
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		s.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエスト形式が不正です",
			"code":  "INVALID_REQUEST_FORMAT",
		})
		return
	}

	// UpdateUserID処理UseCase
	updatedUser, err := s.userUseCase.UpdateUserID(userUpdate.ChangeId, userUpdate.NowId)
	if err != nil {
		s.logger.Error("Failed to update user ID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザーIDの更新に失敗しました",
			"code":  "USER_ID_UPDATE_FAILED",
		})
		return
	}

	//redisでセッション破棄、新IDでセッション作成
	if err := s.sessionManager.UpdateSession(c, userUpdate.ChangeId, userUpdate.NowId); err != nil {
		s.logger.Error("Failed to update session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "セッションの更新に失敗しました",
			"code":  "SESSION_UPDATE_FAILED",
		})
		return
	}

	s.logger.Info("Successfully changed user ID", zap.String("newUserId", userUpdate.ChangeId))
	c.JSON(http.StatusOK, gin.H{
		"message": "ユーザーIDを更新しました",
		"code":    "USER_ID_UPDATED",
		"user":    updatedUser,
	})
}

// 会員情報編集(password)
func (s *SettingController) UpdatePassword(c *gin.Context) {
	var passwordUpdate domain.UserPwChange
	if err := c.ShouldBindJSON(&passwordUpdate); err != nil {
		s.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "リクエスト形式が不正です",
			"code":  "INVALID_REQUEST_FORMAT",
		})
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
		s.logger.Error("Failed to update password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "パスワードの更新に失敗しました",
			"code":  "PASSWORD_UPDATE_FAILED",
		})
		return
	}

	//Redisよりログイン情報セッションを一度消去
	if err := s.sessionManager.DeleteSession(c); err != nil {
		s.logger.Error("Failed to delete session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "セッションの削除に失敗しました",
			"code":  "SESSION_DELETION_FAILED",
		})
		return
	}

	//RedisよりセッションとCookieにUserIdを新しく登録
	if err := s.sessionManager.CreateSession(updatedUser.UserId); err != nil {
		s.logger.Error("Failed to create new session", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "セッションの作成に失敗しました",
			"code":  "SESSION_CREATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "パスワードを更新しました",
		"code":    "PASSWORD_UPDATED",
		"user":    updatedUser,
	})
}
