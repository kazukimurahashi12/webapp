package main

import (
	"log"

	"github.com/kazukimurahashi12/webapp/interface/controller"
	"go.uber.org/zap"
)

var logger *zap.Logger

// ロガーの初期化
func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Panicf("Failed to initialize logger: %v", err)
	}
}

func main() {
	// アプリケーション終了時にロガーをフラッシュ
	if logger != nil {
		_ = logger.Sync()
	}
	log.Println("Start App...")
	controller.GetRouter()
}
