package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/blog"
	"go.uber.org/zap"
)

type EditController struct {
	blogUseCase    blog.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewEditController(blogUseCase blog.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *EditController {
	return &EditController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

func (e *EditController) EditBlog(c *gin.Context) {
	// JSON形式のリクエストボディを構造体にバインドする
	blogPost := domain.BlogPost{}
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		e.logger.Error("Failed to bind JSON request in blog edit", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ブログ編集データの形式が不正です",
			"code":  "INVALID_BLOG_EDIT_FORMAT",
		})
		return
	}

	// セッションからuserIDを取得
	userID, err := e.sessionManager.GetSession(c)
	if err != nil {
		e.logger.Error("Failed to get session", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "セッションが無効です。再度ログインしてください",
			"code":  "SESSION_INVALID",
		})
		return
	}

	// ログインユーザーと編集対象のブログのLoginIDを比較
	if userID != blogPost.LoginId {
		e.logger.Warn("Login user ID does not match blog post's LoginID",
			zap.String("userID", userID),
			zap.String("blogPostLoginID", blogPost.LoginId))
		c.JSON(http.StatusForbidden, gin.H{
			"error": "編集権限がありません",
			"code":  "EDIT_PERMISSION_DENIED",
		})
		return
	}

	// ブログ記事更新処理UseCase
	updatedBlog, err := e.blogUseCase.UpdateBlog(&blogPost)
	if err != nil {
		e.logger.Error("Failed to update blog post", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ブログ記事の更新に失敗しました",
			"code":  "BLOG_UPDATE_FAILED",
		})
		return
	}

	e.logger.Info("Successfully updated blog", zap.Any("blog", updatedBlog))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を更新しました",
		"code":    "BLOG_UPDATED",
		"blog":    updatedBlog,
	})
}
