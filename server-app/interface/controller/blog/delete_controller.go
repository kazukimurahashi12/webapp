package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/blog"
	"go.uber.org/zap"
)

type DeleteController struct {
	blogUseCase    blog.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewDeleteController(blogUseCase blog.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *DeleteController {
	return &DeleteController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// ブログ記事削除
func (d *DeleteController) DeleteBlog(c *gin.Context) {
	// IDをリクエストから取得
	id := c.Param("id")

	// ブログ記事削除処理UseCase
	err := d.blogUseCase.DeleteBlog(id)
	if err != nil {
		d.logger.Error("Failed to delete blog post", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ブログ記事の削除に失敗しました",
			"code":  "BLOG_DELETION_FAILED",
		})
		return
	}

	d.logger.Info("Successfully deleted blog post", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を削除しました",
		"code":    "BLOG_DELETED",
		"blog_id": id,
	})
}
