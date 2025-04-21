package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/infrastructure/di"
	"github.com/kazukimurahashi12/webapp/infrastructure/web/middleware"
)

// APIエンドポイントのルーティング
func GetRouter() *gin.Engine {
	// Ginのルーター作成
	router := gin.Default()

	// CORS設定読み込み
	router.Use(middleware.ConfigureCORS())

	// DIコンテナ作成
	// NewContainer 依存性注入用のコンストラクタ
	container := di.NewContainer()

	// ルーティング設定
	RegisterRoutes(router, container)

	// TODO 必要に応じて改修
	//HTTPSサーバーを起動LSプロトコル使用※ハンドラの登録後に実行登録後に実行
	//第1引数にはポート番号 ":8080" 、第2引数にはTLS証明書のパス、第3引数には秘密鍵のパス
	// router.RunTLS(":8080", "../../certificate/localhost.crt", "../../certificate/localhost.key")

	return router
}
