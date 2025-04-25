package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/infrastructure/web/middleware"
	"github.com/kazukimurahashi12/webapp/interface/dto"
	"github.com/kazukimurahashi12/webapp/interface/mapper"
	"github.com/kazukimurahashi12/webapp/interface/session"
	usecaseUser "github.com/kazukimurahashi12/webapp/usecase/user"
	"go.uber.org/zap"
)

type SettingController struct {
	userUseCase    usecaseUser.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewSettingController(userUseCase usecaseUser.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *SettingController {
	return &SettingController{
		userUseCase:    userUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// 会員情報編集(UserID)
func (s *SettingController) UpdateID(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	var userUpdate dto.UserIdChange
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		s.logger.Error("Failed to bind JSON",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "リクエスト形式が不正です",
			"code":       "INVALID_REQUEST_FORMAT",
			"request_id": requestID,
		})
		return
	}
	// セッションによるログイン認証はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		s.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		s.logger.Error("userID is not a string",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが正しい型ではありません",
			"code":       "USER_ID_TYPE_ERROR",
			"request_id": requestID,
		})
		return
	}

	// UpdateUserID処理UseCase
	updatedUser, err := s.userUseCase.UpdateUserID(userIDStr, userUpdate.NewId)
	if err != nil {
		s.logger.Error("Failed to update user ID",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "ユーザーIDの更新に失敗しました",
			"code":       "USER_ID_UPDATE_FAILED",
			"request_id": requestID,
		})
		return
	}

	//redisでセッション破棄、新IDでセッション作成
	if err := s.sessionManager.UpdateSession(c, userUpdate.NewId); err != nil {
		s.logger.Error("Failed to update session",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "セッションの更新に失敗しました",
			"code":       "SESSION_UPDATE_FAILED",
			"request_id": requestID,
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToUserCreatedResponse(updatedUser)
	s.logger.Info("Successfully changed user ID",
		zap.String("requestID", requestID),
		zap.Any("user", response))

	// ユーザ情報変更完了レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message":    "ユーザーIDを更新しました",
		"code":       "USER_ID_UPDATED",
		"request_id": requestID,
		"user":       response,
	})
}

// 会員情報編集(password)
func (s *SettingController) UpdatePassword(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	var passwordUpdate dto.UserPwChange
	if err := c.ShouldBindJSON(&passwordUpdate); err != nil {
		s.logger.Error("Failed to bind JSON",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "リクエスト形式が不正です",
			"code":       "INVALID_REQUEST_FORMAT",
			"request_id": requestID,
		})
		return
	}
	// セッションによるログイン認証はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		s.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}

	// UpdateUserPassword処理UseCase
	updatedUser, err := s.userUseCase.UpdateUserPassword(
		userID.(string),
		passwordUpdate.NowPassword,
		passwordUpdate.ChangePassword,
	)
	if err != nil {
		s.logger.Error("Failed to update password",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "パスワードの更新に失敗しました",
			"code":       "PASSWORD_UPDATE_FAILED",
			"request_id": requestID,
		})
		return
	}

	//redisでセッション破棄、再度セッション作成
	if err := s.sessionManager.UpdateSession(c, updatedUser.UserID); err != nil {
		s.logger.Error("Failed to update session",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "セッションの更新に失敗しました",
			"code":       "SESSION_UPDATE_FAILED",
			"request_id": requestID,
		})
		return
	}
	// DTOに変換してレスポンス
	response := mapper.ToUserCreatedResponse(updatedUser)
	s.logger.Info("Successfully updated password",
		zap.String("requestID", requestID),
		zap.Any("user", response))

	// ユーザ情報変更完了レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message":    "パスワードを更新しました",
		"code":       "PASSWORD_UPDATED",
		"request_id": requestID,
		"user":       response,
	})
}
