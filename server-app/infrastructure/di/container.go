package di

import (
	"go.uber.org/zap"

	"github.com/kazukimurahashi12/webapp/infrastructure/db"
	"github.com/kazukimurahashi12/webapp/infrastructure/redis"
	"github.com/kazukimurahashi12/webapp/interface/controller/auth"
	authController "github.com/kazukimurahashi12/webapp/interface/controller/auth"
	blogController "github.com/kazukimurahashi12/webapp/interface/controller/blog"
	"github.com/kazukimurahashi12/webapp/interface/controller/common"
	userController "github.com/kazukimurahashi12/webapp/interface/controller/user"
	"github.com/kazukimurahashi12/webapp/interface/session"
	authUseCase "github.com/kazukimurahashi12/webapp/usecase/auth"
	blogUseCase "github.com/kazukimurahashi12/webapp/usecase/blog"
	userUseCase "github.com/kazukimurahashi12/webapp/usecase/user"
)

// Container 依存性注入用の構造体
type Container struct {
	HomeController    *blogController.HomeController
	LoginController   *authController.LoginController
	BlogController    *blogController.BlogController
	SettingController *userController.SettingController
	LogoutController  *auth.LogoutController
	CommonController  *common.CommonController
	SessionManager    session.SessionManager
	logger            *zap.Logger
}

// DI依存性注入用のコンストラクタ
func NewContainer() *Container {
	// Zapロガー初期化
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// セッションストア初期化
	ss := redis.NewRedisSessionStore()

	// DBClient初期化
	dbClient := db.NewDB()

	// Repository初期化
	blogRepo := db.NewBlogRepository(dbClient)
	userRepo := db.NewUserRepository(dbClient)

	// UseCase初期化
	blogUC := blogUseCase.NewBlogUseCase(blogRepo, userRepo)
	authUC := authUseCase.NewAuthUseCase(userRepo)
	userUC := userUseCase.NewUserUseCase(userRepo)

	// Controller初期化
	return &Container{
		HomeController:    blogController.NewHomeController(blogUC, ss, logger),
		LoginController:   authController.NewLoginController(authUC, ss, logger),
		BlogController:    blogController.NewBlogController(blogUC, ss, logger),
		SettingController: userController.NewSettingController(userUC, ss, logger),
		LogoutController:  auth.NewLogoutController(authUC, ss, logger),
		CommonController:  common.NewCommonController(ss, logger),
		SessionManager:    ss,
		logger:            logger,
	}
}
