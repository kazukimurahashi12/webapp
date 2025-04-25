package di

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/kazukimurahashi12/webapp/infrastructure/db"
	"github.com/kazukimurahashi12/webapp/infrastructure/redis"
	"github.com/kazukimurahashi12/webapp/infrastructure/repository"
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
	RegistController  *userController.RegistController
	SettingController *userController.SettingController
	LogoutController  *authController.LogoutController
	CommonController  *common.CommonController
	SessionManager    session.SessionManager
	logger            *zap.Logger
}

// DI依存性注入用のコンストラクタ
func NewContainer() *Container {
	// プロジェクトルートディレクトリを取得
	currentDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory: %w", err))
	}

	// プロジェクトルートディレクトリを移動して取得
	rootDir := filepath.Dir(filepath.Dir(currentDir))
	// PROJECT_ROOT環境変数を設定
	if err := os.Setenv("PROJECT_ROOT", rootDir); err != nil {
		panic(fmt.Errorf("failed to set PROJECT_ROOT environment variable: %w", err))
	}

	// Zapロガー初期化
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// セッションストア初期化
	ss := redis.NewRedisSessionStore()

	// DBManager初期化
	dbManager := db.NewDBManager(logger)
	if err := dbManager.Connect(); err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}
	// DB接続が成功したか確認
	if !dbManager.IsClieintInstance() {
		logger.Error("Database client non-instance")
		os.Exit(1)
	}
	if !dbManager.CheckDBConnection() {
		logger.Error("Database connection failed")
		os.Exit(1)
	}

	// Repository初期化
	blogRepo := repository.NewBlogRepository(dbManager)
	userRepo := repository.NewUserRepository(dbManager)

	// UseCase初期化
	blogUC := blogUseCase.NewBlogUseCase(blogRepo, userRepo)
	authUC := authUseCase.NewAuthUseCase(userRepo)
	userUC := userUseCase.NewUserUseCase(userRepo)

	// Controller初期化
	return &Container{
		HomeController:    blogController.NewHomeController(blogUC, ss, logger),
		LoginController:   authController.NewLoginController(authUC, ss, logger),
		BlogController:    blogController.NewBlogController(blogUC, userUC, ss, logger),
		RegistController:  userController.NewRegistController(userUC, ss, logger),
		SettingController: userController.NewSettingController(userUC, ss, logger),
		LogoutController:  authController.NewLogoutController(authUC, ss, logger),
		CommonController:  common.NewCommonController(ss, logger),
		SessionManager:    ss,
		logger:            logger,
	}
}
