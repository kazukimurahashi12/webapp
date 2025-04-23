package blog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	sessionMocks "github.com/kazukimurahashi12/webapp/interface/session/mocks"
	blogMocks "github.com/kazukimurahashi12/webapp/usecase/blog/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestDeleteController_DeleteBlog(t *testing.T) {
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
		controller := NewDeleteController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.DeleteBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())
	})

	t.Run("BlogNotFound", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodDelete, "/blog/999", nil)
		ctx.Request = req
		ctx.Set("userID", "user123")
		ctx.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		mockBlogUseCase.EXPECT().
			DeleteBlog("999").
			Return(errors.New("not found"))

		logger := zaptest.NewLogger(t)
		controller := NewDeleteController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.DeleteBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	})

	t.Run("UserIDNotFound", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		req := httptest.NewRequest(http.MethodDelete, "/blog/123", nil)
		ctx.Request = req
		// userIDを設定しない

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		logger := zaptest.NewLogger(t)
		controller := NewDeleteController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.DeleteBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})
}
