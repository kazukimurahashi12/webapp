package blog

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/blog"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type HomeController struct {
	blogUseCase    blog.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewHomeController(blogUseCase blog.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *HomeController {
	return &HomeController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// ブログTOP画面表示
func (h *HomeController) GetTop(c *gin.Context) {
	userID, err := h.sessionManager.GetSession(c)
	if err != nil {
		logrus.WithError(err).Error("Failed to get session")
		wrappedErr := errors.Wrap(err, "Failed to get session")
		log.Println(wrappedErr)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "セッションが無効です。再度ログインしてください",
			"code":  "SESSION_INVALID",
		})
		return
	}

	// ブログ記事取得ORM
	blogs, err := h.blogUseCase.GetBlogsByUserID(userID)
	if err != nil {
		h.logger.Error("Failed to get blogs", zap.String("userID", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ブログ記事の取得に失敗しました",
			"code":  "BLOG_FETCH_FAILED",
		})
		return
	}

	h.logger.Debug("Successfully retrieved blogs", zap.String("userID", userID))
	c.JSON(http.StatusOK, gin.H{
		"message": "ブログ記事を取得しました",
		"code":    "BLOG_FETCHED",
		"blogs":   blogs,
		"meta": gin.H{
			"count": len(blogs),
		},
	})
	logrus.Info("@COMPLETE :GetTop")
}

// マイページ表示
func (h *HomeController) GetMypage(c *gin.Context) {
	userID, err := h.sessionManager.GetSession(c)
	if err != nil {
		h.logger.Error("Failed to get session", zap.String("userID", userID), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "セッションが無効です。再度ログインしてください",
			"code":  "SESSION_INVALID",
		})
		return
	}

	// ユーザー情報取得ORM
	user, err := h.blogUseCase.GetUserByID(userID)
	if err != nil {
		h.logger.Error("Failed to get userID", zap.String("userID", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ユーザー情報の取得に失敗しました",
			"code":  "USER_FETCH_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ユーザー情報を取得しました",
		"code":    "USER_FETCHED",
		"user":    user,
	})
	logrus.Info("@COMPLETE :GetMypage")
}
