package main

import (
	"os"

	"github.com/kazukimurahashi12/webapp/interface/controller"
	"go.uber.org/zap"
)

// エラーメッセージ
const (
	ErrDBConnectionFailed = "database connection failed"
	ErrLoggerInitFailed   = "logger initialization failed"
	ErrServerStartFailed  = "failed to start server"
)

var logger *zap.Logger

// logger初期化
func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(ErrLoggerInitFailed)
	}
}

// mainアプリケーションエントリーポイント
func main() {
	defer func() {
		if logger != nil {
			_ = logger.Sync()
		}
	}()

	logger.Info("Starting application...")

	// ルーター初期化
	router := controller.GetRouter()

	// サーバー起動
	logger.Info("Server listening on port 8080")
	if err := router.Run(":8080"); err != nil {
		logger.Error(ErrServerStartFailed, zap.Error(err))
		os.Exit(1)
	}
}
