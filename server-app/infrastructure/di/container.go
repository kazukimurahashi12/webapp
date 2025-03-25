package di

import (
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
}

// Container依存性注入用のコンストラクタ
func NewContainer() *Container {
	//　セッションストア初期化
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
		HomeController:    blogController.NewHomeController(blogUC, ss),
		LoginController:   authController.NewLoginController(authUC, ss),
		BlogController:    blogController.NewBlogController(blogUC, ss),
		SettingController: userController.NewSettingController(userUC, ss),
		LogoutController:  auth.NewLogoutController(authUC, ss),
		CommonController:  common.NewCommonController(ss),
		SessionManager:    ss,
	}
}
