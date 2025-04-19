package blog

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain/blog"
	"github.com/kazukimurahashi12/webapp/domain/user"
	"github.com/kazukimurahashi12/webapp/interface/dto"
	"github.com/kazukimurahashi12/webapp/interface/mapper"
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
	// セッションによるログイン判定はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		b.logger.Error("userID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "userIDが取得できませんでした",
			"code":  "USER_ID_NOT_FOUND",
		})
		return
	}
	// JSON形式のリクエストボディを構造体にバインドする
	req := dto.BlogPost{}
	if err := c.ShouldBindJSON(&req); err != nil {
		b.logger.Error("Failed to bind JSON to struct in blog creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ブログ投稿データの形式が不正です",
			"code":  "INVALID_BLOG_FORMAT",
		})
		return
	}

	// DTO、Entity変換
	entityBlog, err := blog.NewBlog(userID.(string), req.Title, req.Content)
	if err != nil {
		b.logger.Error("Domain validation failed in blog creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_BLOG_ENTITY",
		})
		return
	}

	// ブログ記事登録処理UseCase
	blog, err := b.blogUseCase.NewCreateBlog(entityBlog)
	if err != nil {
		b.logger.Error("Failed to create blog", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ブログ記事の登録に失敗しました",
			"code":  "BLOG_CREATION_FAILED",
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToBlogCreatedResponse(blog)

	b.logger.Info("Successfully created blog", zap.Any("blog", response))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を登録しました",
		"code":    "BLOG_CREATED",
		"blog":    response,
	})
}

// ブログ記事詳細取得
func (b *BlogController) GetBlogView(c *gin.Context) {
	// セッションによるログイン判定はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		b.logger.Error("userID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "userIDが取得できませんでした",
			"code":  "USER_ID_NOT_FOUND",
		})
		return
	}

	// リクエストパラメータからIDを取得
	id := c.Param("id")
	// IDからブログ記事詳細を取得
	blog, err := b.blogUseCase.GetBlogByID(id)
	if err != nil {
		b.logger.Error("Failed to bind JSON to struct in blog creation", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ブログ記事の取得に失敗しました",
			"code":  "BLOG_FETCH_FAILED",
		})
		return
	}

	// 閲覧権限チェック
	if blog.LoginID != userID.(string) {
		b.logger.Warn("Unauthorized blog access attempt", zap.String("blogID", id), zap.String("userID", userID.(string)))
		c.JSON(http.StatusForbidden, gin.H{
			"error": "このブログ記事を閲覧する権限がありません",
			"code":  "BLOG_ACCESS_DENIED",
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToBlogCreatedResponse(blog)
	// 成功時のレスポンス
	b.logger.Info("Successfully fetched blog", zap.Any("blog", response))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を取得しました",
		"code":    "BLOG_FETCHED",
		"blog":    response,
	})
}

// ブログ記事更新
func (b *BlogController) EditBlog(c *gin.Context) {
	// セッションuserIDの取得（ミドルウェアで認証済み）
	userID, exists := c.Get("userID")
	if !exists {
		b.logger.Error("userID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "userIDが取得できませんでした",
			"code":  "USER_ID_NOT_FOUND",
		})
		return
	}

	// リクエストバインド
	req := dto.BlogPost{}
	if err := c.ShouldBindJSON(&req); err != nil {
		b.logger.Error("Failed to bind JSON in blog edit", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ブログ編集データの形式が不正です",
			"code":  "INVALID_BLOG_EDIT_FORMAT",
		})
		return
	}

	// DTO→Entity変換
	entityBlog, err := blog.NewBlog(userID.(string), req.Title, req.Content)
	if err != nil {
		b.logger.Error("Domain validation failed in blog edit", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_BLOG_ENTITY",
		})
		return
	}

	// ブログ更新UseCase
	updatedBlog, err := b.blogUseCase.UpdateBlog(entityBlog)
	if err != nil {
		b.logger.Error("Failed to update blog", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ブログ記事の更新に失敗しました",
			"code":  "BLOG_UPDATE_FAILED",
		})
		return
	}

	response := mapper.ToBlogCreatedResponse(updatedBlog)
	b.logger.Info("Successfully updated blog", zap.Any("blog", response))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を更新しました",
		"code":    "BLOG_UPDATED",
		"blog":    response,
	})
}

// ブログ記事削除
func (b *BlogController) DeleteBlog(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		b.logger.Error("userID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "userIDが取得できませんでした",
			"code":  "USER_ID_NOT_FOUND",
		})
		return
	}

	id := c.Param("id")

	// UseCaseで削除（ユーザーIDによる所有者チェックなども想定）
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
		"message":     "ブログ記事を削除しました",
		"code":        "BLOG_DELETED",
		"blog_id":     id,
		"blog_userID": userID,
	})
}

// 会員情報登録
func (b *BlogController) Regist(c *gin.Context) {
	// セッションから userID を取得（ミドルウェア済み前提）
	userID, exists := c.Get("userID")
	if !exists {
		b.logger.Error("userID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "userIDが取得できませんでした",
			"code":  "USER_ID_NOT_FOUND",
		})
		return
	}

	// JSON形式のリクエストボディを構造体にバインド
	req := dto.FormUser{}
	if err := c.ShouldBindJSON(&req); err != nil {
		b.logger.Error("Failed to bind JSON in user registration", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ユーザー登録データの形式が不正です",
			"code":  "INVALID_USER_REGIST_FORMAT",
		})
		return
	}

	// DTO、Entity変換
	entityUser, err := user.NewUser(userID.(string), req.Password)
	if err != nil {
		b.logger.Error("Domain validation failed in blog creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_BLOG_ENTITY",
		})
		return
	}

	// UseCaseに登録依頼
	user, err := b.blogUseCase.NewCreateUser(entityUser)
	if err != nil {
		b.logger.Error("Failed to register user",
			zap.String("registeredBy", userID.(string)),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザー登録に失敗しました",
			"code":  "USER_REGISTRATION_FAILED",
		})
		return
	}

	b.logger.Info("Successfully registered user",
		zap.String("registeredBy", userID.(string)),
		zap.Any("user", user),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "ユーザー登録が完了しました",
		"code":    "USER_REGISTERED",
		"user":    user,
	})
}
