package blog

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kazukimurahashi12/webapp/domain/blog"
	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
	sessionMocks "github.com/kazukimurahashi12/webapp/interface/session/mocks"
	blogMocks "github.com/kazukimurahashi12/webapp/usecase/blog/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestBlogController_PostBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"title":"test title","content":"test content"}`
		req := httptest.NewRequest(http.MethodPost, "/blog", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req
		ctx.Set("userID", "123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		expectedBlog, _ := blog.NewBlog(123, "test title", "test content")
		mockBlogUseCase.EXPECT().
			FindBlogByAuthorID(uint(123)).
			Return(expectedBlog, nil)

		mockBlogUseCase.EXPECT().
			NewCreateBlog(gomock.Any()).
			Return(expectedBlog, nil)

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		controller.PostBlog(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ブログ記事を登録しました", response["message"])
			assert.Equal(t, "BLOG_CREATED", response["code"])
		}
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"invalid":"data"}`
		req := httptest.NewRequest(http.MethodPost, "/blog", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.PostBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())

		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ブログ投稿データの形式が不正です", response["error"])
			assert.Equal(t, "INVALID_BLOG_FORMAT", response["code"])
		}
	})
}

func TestBlogController_GetBlogView(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const RequestIDKey = "requestID"
	const HeaderXRequestID = "X-Request-ID"

	mockSession := sessionMocks.NewMockSessionManager(ctrl)
	mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)
	logger := zaptest.NewLogger(t)
	controller := NewBlogController(mockBlogUseCase, mockSession, logger)

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/blog/123", nil)
		c.Request = req

		// requestIDを定義
		requestID := "test-request-id"

		// コンテキストに詰める
		reqCtx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(reqCtx)

		// レスポンスヘッダーにも追加
		c.Writer.Header().Set(HeaderXRequestID, requestID)

		c.Set("userID", "123")
		c.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

		expectedBlog := &blog.Blog{
			ID:       123,
			AuthorID: uint(123),
		}

		mockBlogUseCase.EXPECT().
			FindBlogByID(uint(123)).
			Return(expectedBlog, nil)

		controller.GetBlogView(c)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("userID not in context", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/blog/123", nil)
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

		controller.GetBlogView(ctx)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("userID is not a string", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/blog/123", nil)
		ctx.Set("userID", 123)
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

		controller.GetBlogView(ctx)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("invalid blog id format", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/blog/abc", nil)
		ctx.Set("userID", "user123")
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

		controller.GetBlogView(ctx)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("FindBlogByID returns error", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/blog/123", nil)
		ctx.Set("userID", "user123")
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

		mockBlogUseCase.EXPECT().
			FindBlogByID(uint(123)).
			Return(nil, errors.New("not found"))

		controller.GetBlogView(ctx)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run("Unauthorized access", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/blog/123", nil)
		ctx.Set("userID", "user999")
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

		expectedBlog := &blog.Blog{
			ID:       123,
			AuthorID: uint(123),
		}

		mockBlogUseCase.EXPECT().
			FindBlogByID(uint(123)).
			Return(expectedBlog, nil)

		controller.GetBlogView(ctx)
		assert.Equal(t, http.StatusForbidden, recorder.Code)
	})
}

func TestBlogController_EditBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"id":"test id","userID":"user123","title":"test title","content":"test content"}`
		req := httptest.NewRequest(http.MethodPut, "/blog/123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		mockUser := &domainUser.User{
			ID:       123,
			Username: "user123",
		}
		mockBlogUseCase.EXPECT().
			FindBlogByAuthorID("user123").
			Return(mockUser, nil).Times(1)
		expectedBlog, _ := blog.NewBlog(123, "updated title", "updated content")
		mockBlogUseCase.EXPECT().
			UpdateBlog(gomock.Any()).
			Return(expectedBlog, nil)

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.EditBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	})
}

func TestBlogController_DeleteBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodDelete, "/blog/123", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		mockBlogUseCase.EXPECT().
			DeleteBlog(uint(123)).
			Return(nil)

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		controller.DeleteBlog(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("DeleteFailed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodDelete, "/blog/123", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		mockBlogUseCase.EXPECT().
			DeleteBlog(uint(123)).
			Return(errors.New("delete failed"))

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.DeleteBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})
}
