package logout

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/model/db"
	"github.com/kazukimurahashi12/webapp/model/entity"
	"github.com/kazukimurahashi12/webapp/model/redis"
	"github.com/kazukimurahashi12/webapp/service"
)

// ログアウト処理
func DecideLogout(c *gin.Context, redis redis.SessionStore) {
	//構造体をインスタンス化
	logoutUser := entity.FormUser{}
	//JSONデータのリクエストボディを構造体にバインドしてバリデーションを実行
	err := c.ShouldBindJSON(&logoutUser)
	//JSONデータをUser構造体にバインドしてバリデーションを実行
	if err != nil {
		//バリデーションチェックを実行
		err := service.ValidationCheck(c, err)
		if err != nil {

			log.Printf("リクエストJSON形式で構造体にバインドを失敗しました。logoutUser.UserId: %s, err: %v", logoutUser.UserId, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	//DB上の会員情報と照合チェック
	user, err := db.CheckUser(logoutUser.UserId, logoutUser.Password)
	if err != nil {
		log.Printf("ログアウト画面DB上の会員情報と照合に失敗しました。logoutUser.UserId: %s", logoutUser.UserId)
		c.JSON(http.StatusUnauthorized, gin.H{"error in db.CheckUser of decideLogout": err.Error()})
		return
	}

	//Redisよりログイン情報セッションを消去
	err = redis.DeleteSession(c, user.UserId)
	if err != nil {
		log.Println("セッションを消去できませんでした。,err :" + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in decideLogout": err.Error()})
		return
	}

	log.Println("ログアウトに成功しました。Success Logout :id", logoutUser.UserId)
	c.JSON(http.StatusOK, gin.H{"Success Logout": "auth/login"})
}
