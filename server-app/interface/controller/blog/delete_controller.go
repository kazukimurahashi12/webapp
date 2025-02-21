package blog

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/model/db"
)

func GetDeleteBlog(c *gin.Context) {
	//IDをリクエストから取得
	id := c.Param("id")

	blog, err := db.DeleteBlogInfoById(id)
	if err != nil {
		log.Println("ブログ記事の消去に失敗しました。", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Success Deleted Blog :blog.Title %+v", blog.Title)
	c.JSON(http.StatusOK, gin.H{"Deleted blog.Titile": blog.Title})
}
