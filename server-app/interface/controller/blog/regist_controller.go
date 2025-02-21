package blog

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/model/db"
	"github.com/kazukimurahashi12/webapp/model/entity"
	"github.com/kazukimurahashi12/webapp/service"
)

// 新規会員登録(id,password)
func PostRegist(c *gin.Context) {
	//構造体をインスタンス化
	registUser := entity.FormUser{}
	//リクエストをGo構造体にバインド
	err := c.ShouldBindJSON(&registUser)
	//JSONデータをUser構造体にバインドしてバリデーションを実行
	if err != nil {
		//バリデーションチェックを実行
		err := service.ValidationCheck(c, err)
		if err != nil {
			log.Printf("リクエストJSON形式で構造体にバインドを失敗しました。registUser.UserId: %s, err: %v", registUser.UserId, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	//DBに会員情報登録処理
	user, err := db.Signup(registUser.UserId, registUser.Password)
	if err != nil {
		err := errors.New("Error in db.Signup")
		log.Printf("DBに会員情報の登録に失敗しました。registUser.UserId: %s, err: %v", registUser.UserId, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in db.Signup of postRegist": err.Error()})
		return
	}
	//DBに会員情報登録に成功
	log.Printf("Success user in RegisterView from DB :user %+v", user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}
