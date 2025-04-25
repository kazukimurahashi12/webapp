package blog

import (
	"encoding/json"
	"errors"
	"fmt"
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

func TestEditController_EditBlog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Success", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"UserID":"user123","title":"updated title","content":"updated content"}`
		req := httptest.NewRequest(http.MethodPut, "/blog/123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		// uint型に変換
		var testID uint
		if _, err := fmt.Sscanf("user123", "%d", &testID); err != nil {
			t.Logf("Error converting string to uint: %v", err)
		}
		expectedBlog, _ := blog.NewBlog(testID, "updated title", "updated content")
		mockBlogUseCase.EXPECT().
			UpdateBlog(expectedBlog).
			Return(expectedBlog, nil)

		logger := zaptest.NewLogger(t)
		controller := NewEditController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.EditBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, ctx.Writer.Status())

		var response map[string]interface{}
		err := json.Unmarshal(recorder.Body.Bytes(), &response)
		if assert.NoError(t, err) {
			assert.Equal(t, "ブログ記事を更新しました", response["message"])
			assert.Equal(t, "BLOG_UPDATED", response["code"])
		}
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"invalid":"data"}`
		req := httptest.NewRequest(http.MethodPut, "/blog/123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		logger := zaptest.NewLogger(t)
		controller := NewEditController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.EditBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	})

	t.Run("UpdateFailed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"UserID":"user123","title":"updated title","content":"updated content"}`
		req := httptest.NewRequest(http.MethodPut, "/blog/123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req
		ctx.Set("userID", "user123")

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		// モック設定
		// uint型に変換
		var testID uint
		if _, err := fmt.Sscanf("user123", "%d", &testID); err != nil {
			t.Logf("Error converting string to uint: %v", err)
		}
		expectedBlog, _ := blog.NewBlog(testID, "updated title", "updated content")
		mockBlogUseCase.EXPECT().
			UpdateBlog(expectedBlog).
			Return(nil, errors.New("update failed"))

		logger := zaptest.NewLogger(t)
		controller := NewEditController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.EditBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusBadRequest, ctx.Writer.Status())
	})

	t.Run("UserIDNotFound", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		reqBody := `{"UserID":"user123","title":"updated title","content":"updated content"}`
		req := httptest.NewRequest(http.MethodPut, "/blog/123", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req
		// userIDを設定しない

		mockSession := sessionMocks.NewMockSessionManager(ctrl)
		mockBlogUseCase := blogMocks.NewMockUseCase(ctrl)

		logger := zaptest.NewLogger(t)
		controller := NewEditController(mockBlogUseCase, mockSession, logger)

		// 実行
		controller.EditBlog(ctx)

		// 検証
		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status())
	})
}
