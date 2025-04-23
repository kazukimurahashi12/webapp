package blog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kazukimurahashi12/webapp/domain/blog"
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
		reqBody := `{"id":"test id","userID":"user123","title":"test title","content":"test content"}`
		req := httptest.NewRequest(http.MethodPost, "/blog", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		expectedBlog, _ := blog.NewBlog("user123", "test title", "test content")
		mockBlogUseCase.EXPECT().
			NewCreateBlog(gomock.Any()).
			Return(expectedBlog, nil)

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.PostBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())

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

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/blog/123", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		expectedBlog := &blog.Blog{ID: 123, LoginID: "user123"}
		mockBlogUseCase.EXPECT().
			GetBlogByID("123").
			Return(expectedBlog, nil)

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.GetBlogView(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	})

	t.Run("NotFound", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/blog/999", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		mockBlogUseCase.EXPECT().
			GetBlogByID("999").
			Return(nil, errors.New("not found"))

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.GetBlogView(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
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
		expectedBlog, _ := blog.NewBlog("user123", "updated title", "updated content")
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

		// モック設定
		mockBlogUseCase.EXPECT().
			DeleteBlog("123").
			Return(nil)

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.DeleteBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
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
			DeleteBlog("123").
			Return(errors.New("delete failed"))

		logger := zaptest.NewLogger(t)
		controller := NewBlogController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.DeleteBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})
}
