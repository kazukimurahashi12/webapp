package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

// DB接続用初期設定
func NewDB() *DB {
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

	//Db構造体に取得結果代入
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Println("DBの接続に失敗しました。Path:", path)
		connect(dialector, 5)
	}
	return &DB{DB: db}
}

// 接続リトライ処理
func connect(dialector gorm.Dialector, count uint) (*DB, error) {
	db, err := gorm.Open(dialector)
	if err != nil {
		if count > 1 {
			time.Sleep(time.Second * 2)
			count--
			log.Printf("retry... connect to database count:%v\n", count)
			connect(dialector, count)
		}
		log.Println("Failed to connect to database Error:", err.Error())
	}
	if db != nil && err != nil {
		return &DB{DB: db}, nil
	}
	return nil, err
}
