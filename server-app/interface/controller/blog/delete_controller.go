package blog

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/blog"
)

type DeleteController struct {
	blogUseCase    blog.UseCase
	sessionManager session.SessionManager
}

func NewDeleteController(blogUseCase blog.UseCase, sessionManager session.SessionManager) *DeleteController {
	return &DeleteController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
	}
}

// ブログ記事削除
func (d *DeleteController) DeleteBlog(c *gin.Context) {
	// IDをリクエストから取得
	id := c.Param("id")

	// ブログ記事削除処理UseCase
	err := d.blogUseCase.DeleteBlog(id)
	if err != nil {
		log.Printf("ブログ記事の削除に失敗しました。id: %s, error: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Deleted Blog :blog.id %+v", id)
	c.JSON(http.StatusOK, gin.H{"Deleted blog.id": id})
}
