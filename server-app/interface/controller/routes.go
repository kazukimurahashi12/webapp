package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/infrastructure/di"
	"github.com/kazukimurahashi12/webapp/interface/session"
)

// ルーティング設定
func RegisterRoutes(router *gin.Engine, container *di.Container) {
	//共通処理系ルーティング
	router.GET("/", isAuthenticated(container.SessionManager), container.HomeController.GetTop)
	router.GET("/login", container.LoginController.GetLogin)
	router.POST("/login", container.LoginController.PostLogin)

	// Blog系ルーティング
	router.POST("/blog/post", isAuthenticated(container.SessionManager), container.BlogController.PostBlog)
	router.GET("/blog/overview", isAuthenticated(container.SessionManager), container.HomeController.GetMypage)
	router.GET("/blog/overview/post/:id", isAuthenticated(container.SessionManager), container.BlogController.GetBlogView)
	router.POST("/blog/edit", isAuthenticated(container.SessionManager), container.BlogController.EditBlog)
	router.GET("/blog/delete/:id", isAuthenticated(container.SessionManager), container.BlogController.DeleteBlog)

	// User系ルーティング
	router.POST("/update/id", isAuthenticated(container.SessionManager), container.SettingController.UpdateID)
	router.POST("/update/pw", isAuthenticated(container.SessionManager), container.SettingController.UpdatePassword)

	// Auth系ルーティング
	router.POST("/logout", isAuthenticated(container.SessionManager), container.LogoutController.DecideLogout)
	router.POST("/regist", isAuthenticated(container.SessionManager), container.RegistController.Regist)

	// ログイン共通系ルーティング
	router.GET("/api/login-id", isAuthenticated(container.SessionManager), container.CommonController.GetLoginIdBySession)
}

// ログイン中かどうかを判定するミドルウェア
// このハンドラ関数はクライアントのリクエストが処理される前に実行
func isAuthenticated(sessionManager session.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// セッションからユーザーIDを取得
		userID, err := sessionManager.GetSession(c)
		if err != nil {
			log.Println("セッションからIDの取得に失敗しました。", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// セッションにユーザーIDが存在していない場合
		if userID == "" {
			fmt.Println("セッションにユーザーIDが存在していません")
			c.JSON(http.StatusFound, gin.H{"message": "status 302 fail to get session id"})
			//リクエストを中断
			c.Abort()
		}

		// セッションにユーザーIDが存在している場合
		fmt.Println("success get session id")
		// ユーザーIDをコンテキストに保存
		c.Set("userID", userID)
		// 次のミドルウェアまたはハンドラ関数を実行
		c.Next()
	}
}
