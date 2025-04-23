package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/infrastructure/web/middleware"
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
	// コンテクストからリクエストIDを取得
	reqCtx := ctx.Request.Context()
	requestID := middleware.GetRequestID(reqCtx)

	// セッションによるログイン認証はroutes.go_isAuthenticated共通実施しコンテクストから取得
	loginID, exists := ctx.Get("userID")
	if !exists {
		c.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}

	c.logger.Info("Successfully fetched login ID from session",
		zap.String("requestID", requestID),
		zap.String("loginID", loginID.(string)))
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "ログインIDを取得しました",
		"code":       "LOGIN_ID_FETCHED",
		"request_id": requestID,
		"loginID":    loginID.(string),
	})
}
