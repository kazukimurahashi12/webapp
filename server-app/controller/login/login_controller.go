package login

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/model/db"
	"github.com/kazukimurahashi12/webapp/model/entity"
	"github.com/kazukimurahashi12/webapp/model/redis"
	"github.com/kazukimurahashi12/webapp/service"
)

func GetLogin(c *gin.Context, redis redis.SessionStore) {
	user := entity.User{}
	//セッションからuserを取得
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	UserId, err := redis.GetSession(c, cookieKey)
	if err != nil {
		log.Printf("ログイン画面DB上の会員情報のセッションから取得に失敗しました。err: %v", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error in db.GetOneUser of getLogin": err.Error()})
		return
	}

	if UserId != "" {
		var err error
		user, err = db.GetOneUser(UserId)
		if err != nil {
			log.Printf("ログイン画面DB上の会員情報のセッションから取得に失敗しました。err: %v", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error in db.GetOneUser of getLogin": err.Error()})
			return
		}
	}

	//取得成功結果をレスポンス
	log.Printf("Get user in LoginView from DB :user %+v", user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func PostLogin(c *gin.Context, redis redis.SessionStore) {
	//構造体をインスタンス化
	loginUser := entity.FormUser{}
	//JSONデータのリクエストボディを構造体にバインドしてバリデーションを実行
	err := c.ShouldBindJSON(&loginUser)
	//JSONデータをUser構造体にバインドしてバリデーションを実行
	if err != nil {
		//バリデーションチェックを実行
		err := service.ValidationCheck(c, err)
		if err != nil {

			log.Printf("リクエストJSON形式で構造体にバインドを失敗しました。loginUser.UserId: %s, err: %v", loginUser.UserId, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	//DB上の会員情報と照合チェック
	user, err := db.CheckUser(loginUser.UserId, loginUser.Password)
	if err != nil {
		log.Printf("Error in CheckUser ログイン画面DB上の会員情報と照合に失敗しました。loginUser.UserId: %s, err: %v", loginUser.UserId, err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error in db.CheckUser of postLogin": err.Error()})
		return
	}

	//セッションとCookieにUserIdを登録
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	err = redis.NewSession(c, cookieKey, user.UserId)
	if err != nil {
		log.Printf("Error in NewSession ログイン画面DB上のセッションにIDの登録に失敗しました。loginUser.UserId: %s, err: %v", loginUser.UserId, err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error in redis.NewSession of postLogin": err.Error()})
		return
	}

	//取得成功結果をレスポンス
	log.Printf("Success Login :user %+v", user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}
