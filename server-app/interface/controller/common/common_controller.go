package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"go.uber.org/zap"
)

type CommonController struct {
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewCommonController(sessionManager session.SessionManager, logger *zap.Logger) *CommonController {
	return &CommonController{
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// セッションからログインIDを取得するAPI
func (c *CommonController) GetLoginIdBySession(ctx *gin.Context) {
	// セッションからIDを取得
	id, err := c.sessionManager.GetSession(ctx)
	if err != nil {
		c.logger.Error("Failed to get ID from session", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "セッションが無効です。再度ログインしてください",
			"code":  "SESSION_INVALID",
		})
		return
	}

	c.logger.Info("Successfully fetched login ID from session", zap.String("id", id))
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ログインIDを取得しました",
		"code":    "LOGIN_ID_FETCHED",
		"id":      id,
	})
}
