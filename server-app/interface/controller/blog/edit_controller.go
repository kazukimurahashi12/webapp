package blog

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/blog"
)

type EditController struct {
	blogUseCase    blog.UseCase
	sessionManager session.SessionManager
}

func NewEditController(blogUseCase blog.UseCase, sessionManager session.SessionManager) *EditController {
	return &EditController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
	}
}

func (e *EditController) EditBlog(c *gin.Context) {
	// JSON形式のリクエストボディを構造体にバインドする
	blogPost := domain.BlogPost{}
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		log.Printf("Failed to bind JSON request to struct in blog edit. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// セッションからuserIDを取得
	userID, err := e.sessionManager.GetSession(c)
	if err != nil {
		log.Printf("Failed to get session: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// ログインユーザーと編集対象のブログのLoginIDを比較
	if userID != blogPost.LoginId {
		log.Printf("Login user ID does not match blog post's LoginID. userID: %s, blogPost.LoginID: %s", userID, blogPost.LoginId)
		c.JSON(http.StatusForbidden, gin.H{"error": "login user LoginID does not match the target blog's LoginID"})
		return
	}

	// ブログ記事更新処理UseCase
	updatedBlog, err := e.blogUseCase.UpdateBlog(&blogPost)
	if err != nil {
		log.Printf("Failed to update blog post. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Edit Blog :blog %+v", updatedBlog)
	c.JSON(http.StatusOK, gin.H{"blog": updatedBlog})
}
