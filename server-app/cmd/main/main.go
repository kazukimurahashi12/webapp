package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	//起動中ログ出力
	logger.Info("Starting application...")

	// ルーター初期化
	router := controller.GetRouter()

	// ポート設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverAddr := ":" + port

	// HTTPサーバーの作成
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// シグナル処理のためのチャネル
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 別のゴルーチンでサーバーを起動
	go func() {
		logger.Info("Server listening on port " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	// シグナル待機
	<-quit
	logger.Info("Shutting down server...")

	// コンテキストタイムアウト設定
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// サーバーの正常終了
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
