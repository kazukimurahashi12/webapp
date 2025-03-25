package main

import (
	"log"

	"github.com/kazukimurahashi12/webapp/interface/controller"
)

// メイン関数
func main() {
	log.Println("Start App...")
	controller.GetRouter()
}
