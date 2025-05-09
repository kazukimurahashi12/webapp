package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/interface/dto"

	"github.com/kazukimurahashi12/webapp/infrastructure/web/middleware"
	"github.com/kazukimurahashi12/webapp/interface/mapper"
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
	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	// セッションによる認証
	loginID, err := l.sessionManager.GetSession(c)
	if err != nil {
		l.logger.Error("Failed to get session",
			zap.String("requestID", requestID),
			zap.String("loginID", loginID),
			zap.Error(err),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "セッションが無効です。再度ログインしてください",
			"code":       "SESSION_INVALID",
			"request_id": requestID,
		})
		return
	}

	// UseCaseユーザー情報取得
	user, err := l.authUseCase.GetUserByID(loginID)
	if err != nil {
		l.logger.Error("Failed to get user",
			zap.String("requestID", requestID),
			zap.String("userID", loginID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "ユーザー情報の取得に失敗しました",
			"code":       "USER_FETCH_FAILED",
			"request_id": requestID,
		})
		return
	}

	// DTOに変換してレスポンス
	responseUserID := mapper.ToUserIDResponse(user)
	l.logger.Info("Successfully fetched userinfo",
		zap.String("requestID", requestID),
		zap.Any("userID", responseUserID),
	)
	// ユーザ情報取得完了レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message":    "ユーザー情報を取得しました",
		"code":       "USER_FETCHED",
		"request_id": requestID,
		"user":       responseUserID,
	})
}

// ログイン処理
func (l *LoginController) PostLogin(c *gin.Context) {
	// 初期化
	var loginUser dto.FormUser

	// コンテクストからリクエストIDを取得
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	if err := c.ShouldBindJSON(&loginUser); err != nil {
		err := validator.ValidationCheck(c, err)
		if err != nil {
			l.logger.Error("Failed to bind JSON",
				zap.String("requestID", requestID),
				zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"error":      "リクエスト形式が不正です",
				"code":       "INVALID_REQUEST_FORMAT",
				"request_id": requestID,
			})
			return
		}
	}

	// ユーザー認証
	user, err := l.authUseCase.Authenticate(loginUser.UserID, loginUser.Password)
	if err != nil {
		l.logger.Error("Failed to Authorize",
			zap.String("requestID", requestID),
			zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":      "ユーザーIDまたはパスワードが正しくありません",
			"code":       "AUTHENTICATION_FAILED",
			"request_id": requestID,
		})
		return
	}

	// セッション作成
	if err := l.sessionManager.CreateSession(user.Username); err != nil {
		l.logger.Error("Failed to create session",
			zap.String("requestID", requestID),
			zap.String("userID", user.Username),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":      "セッションの作成に失敗しました",
			"code":       "SESSION_CREATION_FAILED",
			"request_id": requestID,
		})
		return
	}

	// DTOに変換してレスポンス
	responseUserID := mapper.ToUserIDResponse(user)
	// ログイン完了レスポンス
	l.logger.Info("Successfully logined",
		zap.String("requestID", requestID),
		zap.Any("userID", responseUserID))
	c.JSON(http.StatusOK, gin.H{
		"message":    "ログインに成功しました",
		"code":       "LOGIN_SUCCESS",
		"request_id": requestID,
		"user":       responseUserID,
	})
}
