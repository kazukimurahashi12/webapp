package blog

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/infrastructure/web/middleware"
	"github.com/kazukimurahashi12/webapp/interface/mapper"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/blog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type HomeController struct {
	blogUseCase    blog.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewHomeController(blogUseCase blog.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *HomeController {
	return &HomeController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// ブログTOP画面表示
func (h *HomeController) GetTop(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	// セッションによるログイン認証はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		h.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}

	// string型変換
	userIDStr, ok := userID.(string)
	if !ok {
		h.logger.Error("userID is not a string",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが正しい型ではありません",
			"code":       "USER_ID_TYPE_ERROR",
			"request_id": requestID,
		})
		return
	}

	// userIDをuintに変換
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid user ID format",
			zap.String("requestID", requestID),
			zap.String("userID", userIDStr),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ユーザーIDの形式が不正です",
			"code":       "INVALID_USER_ID",
			"request_id": requestID,
		})
		return
	}

	// ブログ記事取得ORM
	blogs, err := h.blogUseCase.FindBlogsByAuthorID(uint(userIDUint))
	if err != nil {
		h.logger.Error("Failed to get blogs",
			zap.String("requestID", requestID),
			zap.String("userID", userIDStr),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "ブログ記事の取得に失敗しました",
			"code":       "BLOG_FETCH_FAILED",
			"request_id": requestID,
		})
		return
	}

	h.logger.Debug("Successfully retrieved blogs",
		zap.String("requestID", requestID),
		zap.String("userID", userIDStr))
	c.JSON(http.StatusOK, gin.H{
		"message":    "ブログ記事を取得しました",
		"code":       "BLOG_FETCHED",
		"request_id": requestID,
		"blogs":      blogs,
		"meta": gin.H{
			"count": len(blogs),
		},
	})
	logrus.Info("@COMPLETE :GetTop",
		"requestID", requestID)
}

// マイページ表示
func (h *HomeController) GetMypage(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	// セッションによるログイン認証はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		h.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}

	// string型変換
	userIDStr, ok := userID.(string)
	if !ok {
		h.logger.Error("userID is not a string",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが正しい型ではありません",
			"code":       "USER_ID_TYPE_ERROR",
			"request_id": requestID,
		})
		return
	}

	// userIDをuintに変換
	userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid user ID format",
			zap.String("requestID", requestID),
			zap.String("userID", userIDStr),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ユーザーIDの形式が不正です",
			"code":       "INVALID_USER_ID",
			"request_id": requestID,
		})
		return
	}

	// ユーザー情報取得ORM
	blog, err := h.blogUseCase.FindBlogsByAuthorID(uint(userIDUint))
	if err != nil {
		h.logger.Error("Failed to get userID",
			zap.String("requestID", requestID),
			zap.String("userID", userIDStr),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "ユーザー情報の取得に失敗しました",
			"code":       "USER_FETCH_FAILED",
			"request_id": requestID,
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToBlogsResponse(blog)
	h.logger.Info("Successfully changed blogs",
		zap.String("requestID", requestID),
		zap.Any("blogs", response))
	c.JSON(http.StatusOK, gin.H{
		"message":    "ユーザー情報を取得しました",
		"code":       "USER_FETCHED",
		"request_id": requestID,
		"blogs":      response,
	})
	logrus.Info("@COMPLETE :GetMypage",
		"requestID", requestID)
}
