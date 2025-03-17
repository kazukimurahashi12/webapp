package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func init() {
	//環境変数設定
	//main.goからの相対パス指定
	envErr := godotenv.Load("./build/app/.env")
	if envErr != nil {
		fmt.Println("Error loading .env file", envErr)
	}
}

// APIエンドポイントとクロスオリジンリソース共有（CORS）の設定
func GetRouter() *gin.Engine {
	//ルターを定義
	router := gin.Default()

	// クロスオリジンリソース共有_CORS設定
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://server-app:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Access-Control-Allow-Credentials",
		"Access-Control-Allow-Headers",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"Authorization",
		"Cookie",
	}
	config.AllowCredentials = true
	//クロスオリジンリソース共有を有効化
	router.Use(cors.New(config))

	// RedisSessionManager初期化
	ss := redis.NewRedisSessionStore()

	// DBクライアント初期化（データベースとの接続を確立）
	dbClient := db.NewDB()

	// リポジトリインターフェース依存性注入（データアクセス層の具体実装を作成）
	// Blogリポジトリのインスタンスを生成しDBクライアントを渡す
	blogRepo := db.NewBlogRepository(dbClient)
	// Userリポジトリのインスタンスを生成し、DBクライアントを渡す
	userRepo := db.NewUserRepository(dbClient)

	// ビジネスロジック（ユースケース）のインスタンスを生成
	// BlogUseCaseにBlogリポジトリとUserリポジトリを注入
	blogUC := blogUseCase.NewBlogUseCase(blogRepo, userRepo)
	// AuthUseCaseにUserリポジトリを注入
	authUC := authUseCase.NewAuthUseCase(userRepo)
	// UserUseCaseにUserリポジトリを注入
	userUC := userUseCase.NewUserUseCase(userRepo)

	// HomeControllerのインスタンスを生成しBlogUseCaseを注入
	// HomeControllerプレゼンテーション層
	homeController := blogController.NewHomeController(blogUC, ss)
	// LoginControllerのインスタンスを生成し、AuthUseCaseを注入
	loginController := authController.NewLoginController(authUC, ss)
	// BlogControllerのインスタンスを生成し、AuthUseCaseを注入
	blogController := blogController.NewBlogController(blogUC, ss)
	// SettingControllerのインスタンスを生成し、UserUseCaseを注入
	settingController := userController.NewSettingController(userUC, ss)
	// LogoutControllerのインスタンスを生成し、AuthUseCaseを注入
	logoutController := auth.NewLogoutController(authUC, ss)

	//***ホーム概要画面***
	router.GET("/", isAuthenticated(ss), homeController.GetTop)

	//***ログイン画面***
	router.GET("/login", loginController.GetLogin)
	router.POST("/login", loginController.PostLogin)

	//***ブログ概要画面***
	router.POST("/blog/post", isAuthenticated(ss), blogController.PostBlog)
	//BlogOverview画面
	router.GET("/blog/overview", isAuthenticated(ss), homeController.GetMypage)
	//BlogIDによるView画面
	router.GET("/blog/overview/post/:id", isAuthenticated(ss), blogController.GetBlogView)
	//ブログ記事編集API
	router.POST("/blog/edit", isAuthenticated(ss), blogController.EditBlog)
	//ブログ記事消去API
	router.GET("/blog/delete/:id", isAuthenticated(ss), blogController.DeleteBlog)

	//***ID情報編集画面***
	//ID変更API
	router.POST("/update/id", isAuthenticated(ss), settingController.UpdateID)

	//***PW情報編集画面***
	//PW変更API
	router.POST("/update/pw", isAuthenticated(ss), settingController.UpdatePassword)

	//***ログアウト画面***
	//ログアウト実行API
	router.POST("/logout", isAuthenticated(ss), logoutController.DecideLogout)

	//***会員情報登録画面***
	//登録画面遷移
	router.POST("/regist", isAuthenticated(ss), blogController.Regist)

	// CommonControllerのインスタンスを生成
	commonController := common.NewCommonController(ss)

	//***共通API***
	//セッションからログインIDを取得するAPI
	router.GET("/api/login-id", isAuthenticated(ss), commonController.GetLoginIdBySession)

	//HTTPSサーバーを起動LSプロトコル使用※ハンドラの登録後に実行登録後に実行
	//第1引数にはポート番号 ":8080" 、第2引数にはTLS証明書のパス、第3引数には秘密鍵のパス
	// router.RunTLS(":8080", "../../certificate/localhost.crt", "../../certificate/localhost.key")

	//PORT環境変数で定義
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	} else {
		log.Printf("PORT is %s", port)
	}
	log.Printf("Listening on port %s.", port)
	//HTTPサーバーを起動
	// router.Run(":8080")

	// HTTPサーバーを起動し、エラーログを出力
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("HTTP server failed to start: %v", err)
	}

	return router
}

// ログイン中かどうかを判定するミドルウェア
// このハンドラ関数はクライアントのリクエストが処理される前に実行
func isAuthenticated(sessionManager session.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := sessionManager.GetSession(c)
		if err != nil {
			log.Println("セッションからIDの取得に失敗しました。", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if userID == "" {
			fmt.Println("セッションにユーザーIDが存在していません")
			c.JSON(http.StatusFound, gin.H{"message": "status 302 fail to get session id"})
			c.Abort()
		}
		fmt.Println("success get session id")
		c.Next()
	}
}
