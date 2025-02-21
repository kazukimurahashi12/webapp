package main

import (
	"log"

	"github.com/kazukimurahashi12/webapp/interface/controller"
)

func main() {
	log.Println("Start App...")
	controller.GetRouter()
}
