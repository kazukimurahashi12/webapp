package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/auth"
	"github.com/kazukimurahashi12/webapp/usecase/validator"
)

//#######################################
//ログインコントローラー
//#######################################

type LoginController struct {
	authUseCase    auth.UseCase
	sessionManager session.SessionManager
}

func NewLoginController(authUseCase auth.UseCase, sessionManager session.SessionManager) *LoginController {
	return &LoginController{
		authUseCase:    authUseCase,
		sessionManager: sessionManager,
	}
}

// ログインユーザー情報取得
func (l *LoginController) GetLogin(c *gin.Context) {
	userID, err := l.sessionManager.GetSession(c)
	if err != nil {
		log.Printf("Failed to get session: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// UseCaseユーザー情報取得
	user, err := l.authUseCase.GetUserByID(userID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// ログイン処理
func (l *LoginController) PostLogin(c *gin.Context) {
	var loginUser domain.FormUser
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		err := validator.ValidationCheck(c, err)
		if err != nil {
			log.Printf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// ユーザー認証
	user, err := l.authUseCase.Authenticate(loginUser.UserId, loginUser.Password)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// セッション作成
	if err := l.sessionManager.CreateSession(user.UserId); err != nil {
		log.Printf("Failed to create session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
