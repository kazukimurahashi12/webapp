package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

var Db *gorm.DB

// DB接続用初期設定
func init() {
	//環境変数設定
	//main.goからの相対パス指定
	envErr := godotenv.Load("./build/db/data/.env")
	if envErr != nil {
		fmt.Println("Error loading .env file", envErr)
	}

	// テストケース内で環境変数を設定
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "password")
	os.Setenv("MYSQL_DATABASE", "user_info")
	os.Setenv("MYSQL_LOCAL_HOST", "localhost:3306")

	//環境変数取得
	user := os.Getenv("MYSQL_USER")
	pw := os.Getenv("MYSQL_PASSWORD")
	db_name := os.Getenv("MYSQL_DATABASE")

	var dbHost string
	if os.Getenv("DOCKER_ENV") == "true" {
		// Dockerコンテナ内での接続先を指定
		dbHost = os.Getenv("MYSQL_DOCKER_HOST")
	} else {
		// ローカル環境での接続先を指定
		dbHost = os.Getenv("MYSQL_LOCAL_HOST")
	}
	// PATH設定
	var path string = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, pw, dbHost, db_name)
	dialector := mysql.Open(path)
	var err error
	//Db構造体に取得結果代入
	if Db, err = gorm.Open(dialector); err != nil {
		log.Println("DBの接続に失敗しました。Path:", path)
		connect(dialector, 100)
		// Db = &gorm.DB{} //deploy
	}
	log.Println("DB Connected!!")
}

func connect(dialector gorm.Dialector, count uint) {
	var err error
	if Db, err = gorm.Open(dialector); err != nil {
		if count > 1 {
			time.Sleep(time.Second * 2)
			count--
			log.Printf("retry... connect to database count:%v\n", count)
			connect(dialector, count)
			return
		}
		log.Printf("Failed to connect to database" + err.Error())
	}
}
