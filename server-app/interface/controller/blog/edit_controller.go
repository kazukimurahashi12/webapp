package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain/blog"
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
	// セッションによるログイン判定はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		e.logger.Error("userID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "userIDが取得できませんでした",
			"code":  "USER_ID_NOT_FOUND",
		})
		return
	}
	// JSON形式のリクエストボディを構造体にバインドする
	req := dto.BlogPost{}
	if err := c.ShouldBindJSON(&req); err != nil {
		e.logger.Error("Failed to bind JSON request in blog edit", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ブログ編集データの形式が不正です",
			"code":  "INVALID_BLOG_EDIT_FORMAT",
		})
		return
	}

	// ログインユーザーと編集対象のブログのLoginIDを比較
	if userID.(string) != req.UserID {
		e.logger.Warn("Login user ID does not match blog post's LoginID",
			zap.String("userID", userID.(string)),
			zap.String("blogPostLoginID", req.UserID))
		c.JSON(http.StatusForbidden, gin.H{
			"error": "編集権限がありません",
			"code":  "EDIT_PERMISSION_DENIED",
		})
		return
	}

	// DTO、Entity変換
	entityBlog, err := blog.NewBlog(userID.(string), req.Title, req.Content)
	if err != nil {
		e.logger.Error("Domain validation failed in blog creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_BLOG_ENTITY",
		})
		return
	}

	// ブログ記事更新処理UseCase
	updatedBlog, err := e.blogUseCase.UpdateBlog(entityBlog)
	if err != nil {
		e.logger.Error("Failed to update blog post", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ブログ記事の更新に失敗しました",
			"code":  "BLOG_UPDATE_FAILED",
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToBlogCreatedResponse(updatedBlog)

	// 成功時のレスポンス
	e.logger.Info("Successfully updated blog", zap.Any("blog", response))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を更新しました",
		"code":    "BLOG_UPDATED",
		"blog":    response,
	})
}
