package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/auth"
	"github.com/kazukimurahashi12/webapp/usecase/validator"
	"go.uber.org/zap"
)

//#######################################
// ログインコントローラー
//#######################################

type LoginController struct {
	authUseCase    auth.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewLoginController(authUseCase auth.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *LoginController {
	return &LoginController{
		authUseCase:    authUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// ログインユーザー情報取得
func (l *LoginController) GetLogin(c *gin.Context) {
	loginID, err := l.sessionManager.GetSession(c)
	if err != nil {
		l.logger.Error("Failed to get session", zap.String("loginID", loginID), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "セッションが無効です。再度ログインしてください",
			"code":  "SESSION_INVALID",
		})
		return
	}

	// UseCaseユーザー情報取得
	user, err := l.authUseCase.GetUserByID(loginID)
	if err != nil {
		l.logger.Error("Failed to get user", zap.String("userID", loginID), zap.Error(err))
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
}

// ログイン処理
func (l *LoginController) PostLogin(c *gin.Context) {
	var loginUser domainUser.FormUser
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		err := validator.ValidationCheck(c, err)
		if err != nil {
			l.logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "リクエスト形式が不正です",
				"code":  "INVALID_REQUEST_FORMAT",
			})
			return
		}
	}

	// ユーザー認証
	user, err := l.authUseCase.Authenticate(loginUser.UserID, loginUser.Password)
	if err != nil {
		l.logger.Error("Failed to Authorize", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "ユーザーIDまたはパスワードが正しくありません",
			"code":  "AUTHENTICATION_FAILED",
		})
		return
	}

	// セッション作成
	if err := l.sessionManager.CreateSession(user.UserID); err != nil {
		l.logger.Error("Failed to create session", zap.String("userID", user.UserID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "セッションの作成に失敗しました",
			"code":  "SESSION_CREATION_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ログインに成功しました",
		"code":    "LOGIN_SUCCESS",
		"user":    user,
	})
}
