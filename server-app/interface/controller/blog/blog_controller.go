package blog

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
	"github.com/kazukimurahashi12/webapp/infrastructure/web/middleware"
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
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	// セッションによるログイン認証はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		b.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		b.logger.Error("userID is not a string",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが正しい型ではありません",
			"code":       "USER_ID_TYPE_ERROR",
			"request_id": requestID,
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
	// 文字列のuserIDをuintに変換
	authorID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		b.logger.Error("Failed to parse userID",
			zap.String("userID", userIDStr),
			zap.Error(err),
			zap.String("requestID", requestID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ユーザーIDの形式が不正です",
			"code":       "INVALID_USER_ID_FORMAT",
			"request_id": requestID,
		})
		return
	}

	// 著者IDからブログを取得
	blog, err := b.blogUseCase.FindBlogByAuthorID(uint(authorID))
	if err != nil {
		b.logger.Error("Failed to find blog by authorID",
			zap.Uint("authorID", uint(authorID)),
			zap.Error(err),
			zap.String("requestID", requestID),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"error":      "指定された著者のブログが存在しません",
			"code":       "BLOG_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}

	// DTO、Entity変換
	entityBlog, err := domainBlog.NewBlog(blog.AuthorID, req.Title, req.Content)
	if err != nil {
		b.logger.Error("Domain validation failed in blog creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_BLOG_ENTITY",
		})
		return
	}

	// ブログ記事登録処理UseCase
	createdBlog, err := b.blogUseCase.NewCreateBlog(entityBlog)
	if err != nil {
		b.logger.Error("Failed to create blog", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ブログ記事の登録に失敗しました",
			"code":  "BLOG_CREATION_FAILED",
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToBlogCreatedResponse(createdBlog)

	b.logger.Info("Successfully created blog", zap.Any("blog", response))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を登録しました",
		"code":    "BLOG_CREATED",
		"blog":    response,
	})
}

// ブログ記事詳細取得
func (b *BlogController) GetBlogView(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	// セッションによるログイン認証はroutes.go_isAuthenticated共通実施しコンテクストから取得
	userID, exists := c.Get("userID")
	if !exists {
		b.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		b.logger.Error("userID is not a string",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが正しい型ではありません",
			"code":       "USER_ID_TYPE_ERROR",
			"request_id": requestID,
		})
		return
	}

	// リクエストパラメータからIDを取得
	idStr := c.Param("id")
	// uint型に変換
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		b.logger.Error("Invalid blog ID format",
			zap.String("requestID", requestID),
			zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ブログIDの形式が不正です",
			"code":       "INVALID_BLOG_ID",
			"request_id": requestID,
		})
		return
	}

	// IDからブログ記事詳細を取得
	blog, err := b.blogUseCase.FindBlogByID(id)
	if err != nil {
		b.logger.Error("Failed to get blog by ID",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "ブログ記事の取得に失敗しました",
			"code":       "BLOG_FETCH_FAILED",
			"request_id": requestID,
		})
		return
	}

	// 閲覧権限チェック
	if blog.Author.Username != userIDStr {
		b.logger.Warn("Unauthorized blog access attempt",
			zap.String("requestID", requestID),
			zap.Uint("blogID", id),
			zap.String("userID", userIDStr))
		c.JSON(http.StatusForbidden, gin.H{
			"error":      "このブログ記事を閲覧する権限がありません",
			"code":       "BLOG_ACCESS_DENIED",
			"request_id": requestID,
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToBlogCreatedResponse(blog)
	// 成功時のレスポンス
	b.logger.Info("Successfully fetched blog",
		zap.String("requestID", requestID),
		zap.Any("blog", response))
	c.JSON(http.StatusOK, gin.H{
		"message":    "ブログ記事を取得しました",
		"code":       "BLOG_FETCHED",
		"request_id": requestID,
		"blog":       response,
	})
}

// ブログ記事更新
func (b *BlogController) EditBlog(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	// セッションuserIDの取得（ミドルウェアで認証済み）
	userID, exists := c.Get("userID")
	if !exists {
		b.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		b.logger.Error("userID is not a string",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが正しい型ではありません",
			"code":       "USER_ID_TYPE_ERROR",
			"request_id": requestID,
		})
		return
	}

	// リクエストバインド
	req := dto.BlogPost{}
	if err := c.ShouldBindJSON(&req); err != nil {
		b.logger.Error("Failed to bind JSON in blog edit",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ブログ編集データの形式が不正です",
			"code":       "INVALID_BLOG_EDIT_FORMAT",
			"request_id": requestID,
		})
		return
	}

	// 文字列のuserIDをuintに変換
	authorID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		b.logger.Error("Failed to parse userID",
			zap.String("userID", userIDStr),
			zap.Error(err),
			zap.String("requestID", requestID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ユーザーIDの形式が不正です",
			"code":       "INVALID_USER_ID_FORMAT",
			"request_id": requestID,
		})
		return
	}

	// DTO→Entity変換
	entityBlog, err := domainBlog.NewBlog(uint(authorID), req.Title, req.Content)
	if err != nil {
		b.logger.Error("Domain validation failed in blog edit",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      err.Error(),
			"code":       "INVALID_BLOG_ENTITY",
			"request_id": requestID,
		})
		return
	}

	// ブログ更新UseCase
	updatedBlog, err := b.blogUseCase.UpdateBlog(entityBlog)
	if err != nil {
		b.logger.Error("Failed to update blog",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "ブログ記事の更新に失敗しました",
			"code":       "BLOG_UPDATE_FAILED",
			"request_id": requestID,
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToBlogCreatedResponse(updatedBlog)
	b.logger.Info("Successfully updated blog",
		zap.String("requestID", requestID),
		zap.Any("blog", response))
	c.JSON(http.StatusOK, gin.H{
		"message":    "ブログ記事を更新しました",
		"code":       "BLOG_UPDATED",
		"request_id": requestID,
		"blog":       response,
	})
}

// ブログ記事削除
func (b *BlogController) DeleteBlog(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	userID, exists := c.Get("userID")
	if !exists {
		b.logger.Error("userID not found in context",
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "userIDが取得できませんでした",
			"code":       "USER_ID_NOT_FOUND",
			"request_id": requestID,
		})
		return
	}

	idStr := c.Param("id")
	// uint型に変換
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		b.logger.Error("Invalid blog ID format",
			zap.String("requestID", requestID),
			zap.String("id", idStr))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ブログIDの形式が不正です",
			"code":       "INVALID_BLOG_ID",
			"request_id": requestID,
		})
		return
	}

	// UseCaseで削除（ユーザーIDによる所有者チェックなども想定）
	err := b.blogUseCase.DeleteBlog(id)
	if err != nil {
		b.logger.Error("Failed to delete blog",
			zap.String("requestID", requestID),
			zap.String("id", fmt.Sprintf("%d", id)),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "ブログ記事の削除に失敗しました",
			"code":       "BLOG_DELETION_FAILED",
			"request_id": requestID,
		})
		return
	}

	b.logger.Info("Successfully deleted blog",
		zap.String("requestID", requestID),
		zap.String("id", fmt.Sprintf("%d", id)))
	c.JSON(http.StatusOK, gin.H{
		"message":     "ブログ記事を削除しました",
		"code":        "BLOG_DELETED",
		"request_id":  requestID,
		"blog_id":     id,
		"blog_userID": userID,
	})
}
