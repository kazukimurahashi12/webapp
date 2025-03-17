package blog

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/infrastructure/redis"
	"github.com/pkg/errors"
)

func PostEditBlog(c *gin.Context, redis redis.RedisSessionStore) {
	// JSON形式のリクエストボディを構造体にバインドする
	blogPost := domain.BlogPost{}
	if err := c.ShouldBindJSON(&blogPost); err != nil {
		log.Printf("ブログ編集画面リクエストJSON形式で構造体にバインドを失敗しました。" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in c.ShouldBindJSON": err.Error()})
		return
	}

	//セッションからloginIDを取得
	userID, err := l.sessionManager.GetSession(c)
	if err != nil {
		log.Printf("Failed to get session: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
	if err := blogRepo.Update(&blog); err != nil {
		log.Println("error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in db.Edit": err.Error()})
		return
	}

	log.Printf("Success Edit Blog :blog %+v", blog)
	c.JSON(http.StatusOK, gin.H{"blog": blog})
}
