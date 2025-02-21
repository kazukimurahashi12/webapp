package blog

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/kazukimurahashi12/webapp/model/redis"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/model/db"
	"github.com/kazukimurahashi12/webapp/model/entity"
)

func PostEditBlog(c *gin.Context, redis redis.SessionStore) {
	// JSON形式のリクエストボディを構造体にバインドする
	blogPost := entity.BlogPost{}
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		log.Printf("ブログ編集画面リクエストJSON形式で構造体にバインドを失敗しました。" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in c.ShouldBindJSON": err.Error()})
		return
	}

	//セッションからloginIDを取得
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	id, err := redis.GetSession(c, cookieKey)
	if err != nil {
		log.Println("セッションからIDの取得に失敗しました。", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// ログインユーザーと編集対象のブログのLoginIDを比較
	if id != blogPost.LoginID {
		err := errors.New("ログインユーザーと編集対象のブログのLoginIDが一致しません。")
		log.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error in blogPost.LoginID": err.Error()})
		return
	}

	//DBにブログ記事内容を登録
	blog, err := db.Edit(blogPost.ID, blogPost.LoginID, blogPost.Title, blogPost.Content)
	if err != nil {
		log.Println("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in db.Edit": err.Error()})
		return
	}

	log.Printf("Success Edit Blog :blog %+v", blog)
	c.JSON(http.StatusOK, gin.H{"blog": blog})
}
