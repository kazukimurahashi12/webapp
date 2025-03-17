package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/user"
	"github.com/kazukimurahashi12/webapp/usecase/validator"
)

type RegistController struct {
	userUseCase    user.UseCase
	sessionManager session.SessionManager
}

func NewRegistController(userUseCase user.UseCase, sessionManager session.SessionManager) *RegistController {
	return &RegistController{
		userUseCase:    userUseCase,
		sessionManager: sessionManager,
	}
}

// 新規会員登録
func (r *RegistController) Regist(c *gin.Context) {
	// JSON形式のリクエストボディを構造体にバインドする
	registUser := domain.FormUser{}
	if err := c.ShouldBindJSON(&registUser); err != nil {
		// バリデーションチェックを実行
		err := validator.ValidationCheck(c, err)
		if err != nil {
			log.Printf("Failed to bind JSON request to struct. userId: %s, error: %v", registUser.UserId, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// 会員情報登録処理UseCase
	createdUser, err := r.userUseCase.CreateUser(&registUser)
	if err != nil {
		log.Printf("Failed to register user. userId: %s, error: %v", registUser.UserId, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Regist User :user %+v", createdUser)
	c.JSON(http.StatusOK, gin.H{"user": createdUser})
}
