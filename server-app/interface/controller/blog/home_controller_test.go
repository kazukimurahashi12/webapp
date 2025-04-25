package blog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kazukimurahashi12/webapp/domain/blog"
	"github.com/kazukimurahashi12/webapp/domain/user"
	sessionMocks "github.com/kazukimurahashi12/webapp/interface/session/mocks"
	blogMocks "github.com/kazukimurahashi12/webapp/usecase/blog/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestHomeController_GetTop(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/blog/top", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		expectedBlogs := []blog.Blog{
			{ID: 1, Title: "Test Blog 1"},
			{ID: 2, Title: "Test Blog 2"},
		}
		mockBlogUseCase.EXPECT().
			FindBlogsByUserID("user123").
			Return(expectedBlogs, nil)

		logger := zaptest.NewLogger(t)
		controller := NewHomeController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.GetTop(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())

		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ブログ記事を取得しました", response["message"])
			assert.Equal(t, "BLOG_FETCHED", response["code"])
			assert.Len(t, response["blogs"], 2)
		}
	})

	t.Run("BlogFetchFailed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/blog/top", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		mockBlogUseCase.EXPECT().
			FindBlogsByUserID("user123").
			Return(nil, errors.New("fetch failed"))

		logger := zaptest.NewLogger(t)
		controller := NewHomeController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.GetTop(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})

	t.Run("UserIDNotFound", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/blog/top", nil)
		ctx.Request = req
		// userIDを設定しない

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		logger := zaptest.NewLogger(t)
		controller := NewHomeController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.GetTop(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})
}

func TestHomeController_GetMypage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/blog/mypage", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		expectedUser := &user.User{UserID: "user123"}
		mockBlogUseCase.EXPECT().
			FindBlogsByUserID("user123").
			Return(expectedUser, nil)

		logger := zaptest.NewLogger(t)
		controller := NewHomeController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.GetMypage(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())

		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ユーザー情報を取得しました", response["message"])
			assert.Equal(t, "USER_FETCHED", response["code"])
		}
	})

	t.Run("UserFetchFailed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodGet, "/blog/mypage", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		mockBlogUseCase.EXPECT().
			FindBlogsByUserID("user123").
			Return(nil, errors.New("fetch failed"))

		logger := zaptest.NewLogger(t)
		controller := NewHomeController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.GetMypage(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})
}
