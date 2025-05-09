package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
	"github.com/kazukimurahashi12/webapp/infrastructure/web/middleware"
	"github.com/kazukimurahashi12/webapp/interface/dto"
	"github.com/kazukimurahashi12/webapp/interface/mapper"
	"github.com/kazukimurahashi12/webapp/interface/session"
	"github.com/kazukimurahashi12/webapp/usecase/user"
	"github.com/kazukimurahashi12/webapp/usecase/validator"
	"go.uber.org/zap"
)

type RegistController struct {
	userUseCase    user.UseCase
	sessionManager session.SessionManager
	logger         *zap.Logger
}

func NewRegistController(userUseCase user.UseCase, sessionManager session.SessionManager, logger *zap.Logger) *RegistController {
	return &RegistController{
		userUseCase:    userUseCase,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// 新規会員登録
func (r *RegistController) Regist(c *gin.Context) {
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	// JSON形式のリクエストボディを構造体にバインドする
	dtoUser := dto.FormUser{}
	if err := c.ShouldBindJSON(&dtoUser); err != nil {
		// バリデーションチェックを実行
		err := validator.ValidationCheck(c, err)
		if err != nil {
			r.logger.Error("Failed to bind JSON request",
				zap.String("requestID", requestID),
				zap.String("username", dtoUser.UserID),
				zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error":      "リクエスト形式が不正です",
				"code":       "INVALID_REQUEST_FORMAT",
				"request_id": requestID,
			})
			return
		}
	}

	// DTO、Entity変換
	entityUser, err := domainUser.NewUser(dtoUser.UserID, dtoUser.Password)
	if err != nil {
		r.logger.Error("Domain validation failed in blog creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_BLOG_ENTITY",
		})
		return
	}

	// 会員情報登録処理UseCase
	createdUser, err := r.userUseCase.CreateUser(entityUser.Username, entityUser.Password)
	if err != nil {
		r.logger.Error("Failed to register user",
			zap.String("requestID", requestID),
			zap.String("username", entityUser.Username),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "ユーザー登録に失敗しました",
			"code":       "USER_REGISTATION_FAILED",
			"request_id": requestID,
		})
		return
	}

	// DTOに変換してレスポンス
	response := mapper.ToUserCreatedResponse(createdUser)
	r.logger.Info("Successfully registered user",
		zap.String("requestID", requestID),
		zap.Any("user", response))

	// ユーザー登録完了レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message":    "ユーザー登録が完了しました",
		"code":       "USER_REGISTERED",
		"request_id": requestID,
		"user":       response,
	})
}
