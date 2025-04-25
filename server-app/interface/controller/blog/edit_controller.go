package blog

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain/blog"
	"github.com/kazukimurahashi12/webapp/infrastructure/web/middleware"
	"github.com/kazukimurahashi12/webapp/interface/dto"
	"github.com/kazukimurahashi12/webapp/interface/mapper"
	"github.com/kazukimurahashi12/webapp/interface/session"
	usecaseBlog "github.com/kazukimurahashi12/webapp/usecase/blog"
	"go.uber.org/zap"
)

type EditController struct {
	blogUseCase    usecaseBlog.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewEditController(blogUseCase usecaseBlog.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *EditController {
	return &EditController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// ブログ記事編集
func (e *EditController) EditBlog(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	// セッションによるログイン認証はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		e.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}
	// JSON形式のリクエストボディを構造体にバインドする
	req := dto.BlogPost{}
	if err := c.ShouldBindJSON(&req); err != nil {
		e.logger.Error("Failed to bind JSON request in blog edit",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ブログ編集データの形式が不正です",
			"code":       "INVALID_BLOG_EDIT_FORMAT",
			"request_id": requestID,
		})
		return
	}

	// string型変換
	userIDStr, ok := userID.(string)
	if !ok {
		e.logger.Error("userID is not a string",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが正しい型ではありません",
			"code":       "USER_ID_TYPE_ERROR",
			"request_id": requestID,
		})
		return
	}

	// ログインユーザーと編集対象のブログのLoginIDを比較
	if userIDStr != req.UserID {
		e.logger.Warn("Login user ID does not match blog post's LoginID",
			zap.String("requestID", requestID),
			zap.String("userID", userIDStr),
			zap.String("blogPostLoginID", req.UserID))
		c.JSON(http.StatusForbidden, gin.H{
			"error":      "編集権限がありません",
			"code":       "EDIT_PERMISSION_DENIED",
			"request_id": requestID,
		})
		return
	}

	// uint型に変換
	var id uint
	if _, err := fmt.Sscanf(req.ID, "%d", &id); err != nil {
		e.logger.Error("Invalid blog ID format",
			zap.String("requestID", requestID),
			zap.String("id", req.ID))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ブログIDの形式が不正です",
			"code":       "INVALID_BLOG_ID",
			"request_id": requestID,
		})
		return
	}
	// DTO、Entity変換
	entityBlog, err := blog.NewBlog(id, req.Title, req.Content)
	if err != nil {
		e.logger.Error("Domain validation failed in blog creation",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      err.Error(),
			"code":       "INVALID_BLOG_ENTITY",
			"request_id": requestID,
		})
		return
	}

	// ブログ記事更新処理UseCase
	updatedBlog, err := e.blogUseCase.UpdateBlog(entityBlog)
	if err != nil {
		e.logger.Error("Failed to update blog post",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ブログ記事の更新に失敗しました",
			"code":       "BLOG_UPDATE_FAILED",
			"request_id": requestID,
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToBlogCreatedResponse(updatedBlog)

	// 成功時のレスポンス
	e.logger.Info("Successfully updated blog",
		zap.String("requestID", requestID),
		zap.Any("blog", response))
	c.JSON(http.StatusOK, gin.H{
		"message":    "ブログ記事を更新しました",
		"code":       "BLOG_UPDATED",
		"request_id": requestID,
		"blog":       response,
	})
}
