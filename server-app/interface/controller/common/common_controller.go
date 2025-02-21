package common

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/model/redis"
)

// セッションからログインIDを取得するAPI
func GetLoginIdBySession(c *gin.Context, redis redis.SessionStore) {
	//セッションからuserを取得
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	//セッションから取得
	id, err := redis.GetSession(c, cookieKey)
	if err != nil {
		log.Println("セッションからIDの取得に失敗しました。", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Get LoginId bySession :id %+v", id)
	c.JSON(http.StatusOK, gin.H{"id": id})
}
