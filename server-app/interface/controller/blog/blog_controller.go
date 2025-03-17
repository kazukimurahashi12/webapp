package blog

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/session"
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
	// JSON形式のリクエストボディを構造体にバインドする
	blogPost := domain.BlogPost{}
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		log.Printf("ブログ記事作成画面でJSON形式構造体にバインドを失敗しました。" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ブログ記事登録処理UseCase
	blog, err := b.blogUseCase.NewCreateBlog(&blogPost)
	if err != nil {
		log.Println("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// DBにブログ記事登録に成功
	log.Printf("Success Post Blog :blog %+v", blog)
	c.JSON(http.StatusOK, gin.H{"blog": blog})
}

// ブログ記事詳細取得
func (b *BlogController) GetBlogView(c *gin.Context) {
	id := c.Param("id")

	blog, err := b.blogUseCase.GetBlogByID(id)
	if err != nil {
		log.Printf("ブログ記事の取得に失敗しました。id: %s, error: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Get Blog View :blog %+v", blog)
	c.JSON(http.StatusOK, gin.H{"blog": blog})
}

// ブログ記事編集
func (b *BlogController) EditBlog(c *gin.Context) {
	// JSON形式のリクエストボディを構造体にバインドする
	blogPost := domain.BlogPost{}
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		log.Printf("ブログ記事編集画面でJSON形式構造体にバインドを失敗しました。error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ブログ記事更新処理UseCase
	updatedBlog, err := b.blogUseCase.UpdateBlog(&blogPost)
	if err != nil {
		log.Printf("ブログ記事の更新に失敗しました。error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Edit Blog :blog %+v", updatedBlog)
	c.JSON(http.StatusOK, gin.H{"blog": updatedBlog})
}

// ブログ記事削除
func (b *BlogController) DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	err := b.blogUseCase.DeleteBlog(id)
	if err != nil {
		log.Printf("ブログ記事の削除に失敗しました。id: %s, error: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Deleted Blog :blog.id %+v", id)
	c.JSON(http.StatusOK, gin.H{"Deleted blog.id": id})
}

// 会員情報登録
func (b *BlogController) Regist(c *gin.Context) {
	// JSON形式のリクエストボディを構造体にバインドする
	registUser := domain.FormUser{}
	if err := c.ShouldBindJSON(&registUser); err != nil {
		log.Printf("会員情報登録画面でJSON形式構造体にバインドを失敗しました。error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 会員情報登録処理UseCase
	user, err := b.blogUseCase.NewCreateUser(&registUser)
	if err != nil {
		log.Printf("会員情報の登録に失敗しました。error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Regist User :user %+v", user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}
