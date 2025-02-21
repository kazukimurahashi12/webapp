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
}

func NewHomeController(blogUseCase blog.UseCase, sessionManager session.SessionManager) *HomeController {
	return &HomeController{
		blogUseCase:    blogUseCase,
		sessionManager: sessionManager,
	}
}

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

func (h *HomeController) GetTop(c *gin.Context) {
	userID, err := h.sessionManager.GetSession(c)
	if err != nil {
		logrus.WithError(err).Error("Failed to get session")
		wrappedErr := errors.Wrap(err, "Failed to get session")
		log.Println(wrappedErr) // スタックトレース付きエラーログ

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authentication failed. Please login again.",
			"code":  "AUTH_ERROR",
		})
		return
	}

	blogs, err := h.blogUseCase.GetBlogsByUserID(userID)
	if err != nil {
		logger.Error("Failed to get blogs", zap.String("userID", userID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve blog data. Please try again later.",
			"code":  "BLOG_FETCH_ERROR",
		})
		return
	}

	logger.Debug("Successfully retrieved blogs", zap.String("userID", userID))
	c.JSON(http.StatusOK, gin.H{
		"blogs": blogs,
		"meta": gin.H{
			"count": len(blogs),
		},
	})
	logrus.Info("@COMPLETE :GetTop")
}

func (h *HomeController) GetMypage(c *gin.Context) {
	userID, err := h.sessionManager.GetSession(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := h.blogUseCase.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
	logrus.Info("@COMPLETE :GetMypage")
	c.Next()
}
