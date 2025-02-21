package blog

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/model/db"
	"github.com/kazukimurahashi12/webapp/model/entity"
	"github.com/kazukimurahashi12/webapp/usecase/blog"
)

type BlogController struct {
	blogUseCase    blog.UseCase
	sessionManager session.SessionManager
}

func NewBlogController(blogUseCase blog.UseCase, sessionManager session.SessionManager) *BlogController {
	return &BlogController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
	}
}

// blog記事登録
func (b *BlogController) PostBlog(c *gin.Context) {

}

func PostBlogOLD(c *gin.Context) {
	// JSON形式のリクエストボディを構造体にバインドする
	blogPost := entity.BlogPost{}
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		log.Printf("ブログ記事作成画面でJSON形式構造体にバインドを失敗しました。" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//DBにブログ記事内容を登録
	blog, err := db.Create(blogPost.LoginID, blogPost.Title, blogPost.Content)
	if err != nil {
		log.Println("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Post Blog :blog %+v", blog)
	c.JSON(http.StatusOK, gin.H{"blog": blog})
}
