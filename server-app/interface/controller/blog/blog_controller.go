package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
	"github.com/kazukimurahashi12/webapp/interface/session"
	usecaseBlog "github.com/kazukimurahashi12/webapp/usecase/blog"
	"go.uber.org/zap"
)

type BlogController struct {
	blogUseCase    usecaseBlog.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewBlogController(blogUseCase usecaseBlog.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *BlogController {
	return &BlogController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// blog記事登録
func (b *BlogController) PostBlog(c *gin.Context) {
	// JSON形式のリクエストボディを構造体にバインドする
	blogPost := domainBlog.BlogPost{}
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		b.logger.Error("Failed to bind JSON to struct in blog creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ブログ投稿データの形式が不正です",
			"code":  "INVALID_BLOG_FORMAT",
		})
		return
	}

	// ブログ記事登録処理UseCase
	blog, err := b.blogUseCase.NewCreateBlog(&blogPost)
	if err != nil {
		b.logger.Error("Failed to create blog", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ブログ記事の登録に失敗しました",
			"code":  "BLOG_CREATION_FAILED",
		})
		return
	}

	b.logger.Info("Successfully created blog", zap.Any("blog", blog))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を登録しました",
		"code":    "BLOG_CREATED",
		"blog":    blog,
	})
}

// ブログ記事詳細取得
func (b *BlogController) GetBlogView(c *gin.Context) {
	id := c.Param("id")

	blog, err := b.blogUseCase.GetBlogByID(id)
	if err != nil {
		b.logger.Error("Failed to bind JSON to struct in blog creation", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ブログ記事の取得に失敗しました",
			"code":  "BLOG_FETCH_FAILED",
		})
		return
	}

	b.logger.Info("Successfully fetched blog", zap.Any("blog", blog))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を取得しました",
		"code":    "BLOG_FETCHED",
		"blog":    blog,
	})
}

// ブログ記事編集
func (b *BlogController) EditBlog(c *gin.Context) {
	// JSON形式のリクエストボディを構造体にバインドする
	blogPost := domainBlog.BlogPost{}
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		b.logger.Error("Failed to bind JSON in blog edit", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ブログ編集データの形式が不正です",
			"code":  "INVALID_BLOG_EDIT_FORMAT",
		})
		return
	}

	// ブログ記事更新処理UseCase
	updatedBlog, err := b.blogUseCase.UpdateBlog(&blogPost)
	if err != nil {
		b.logger.Error("Failed to update blog", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ブログ記事の更新に失敗しました",
			"code":  "BLOG_UPDATE_FAILED",
		})
		return
	}

	b.logger.Info("Successfully updated blog", zap.Any("blog", updatedBlog))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を更新しました",
		"code":    "BLOG_UPDATED",
		"blog":    updatedBlog,
	})
}

// ブログ記事削除
func (b *BlogController) DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	err := b.blogUseCase.DeleteBlog(id)
	if err != nil {
		b.logger.Error("Failed to delete blog", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ブログ記事の削除に失敗しました",
			"code":  "BLOG_DELETION_FAILED",
		})
		return
	}

	b.logger.Info("Successfully deleted blog", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を削除しました",
		"code":    "BLOG_DELETED",
		"blog_id": id,
	})
}

// 会員情報登録
func (b *BlogController) Regist(c *gin.Context) {
	// JSON形式のリクエストボディを構造体にバインドする
	registUser := domainUser.FormUser{}
	if err := c.ShouldBindJSON(&registUser); err != nil {
		b.logger.Error("Failed to bind JSON in user registration", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ユーザー登録データの形式が不正です",
			"code":  "INVALID_USER_REGIST_FORMAT",
		})
		return
	}

	// 会員情報登録処理UseCase
	user, err := b.blogUseCase.NewCreateUser(&registUser)
	if err != nil {
		b.logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザー登録に失敗しました",
			"code":  "USER_REGISTRATION_FAILED",
		})
		return
	}

	b.logger.Info("Successfully registered user", zap.Any("user", user))
	c.JSON(http.StatusOK, gin.H{
		"message": "ユーザー登録が完了しました",
		"code":    "USER_REGISTERED",
		"user":    user,
	})
}
