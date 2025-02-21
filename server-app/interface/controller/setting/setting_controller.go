package setting

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

// 会員情報編集(id)
func PostUpdateId(c *gin.Context, redis redis.SessionStore) {
	user := entity.UserIdChange{}
	//リクエストをGo構造体にバインド
	err := c.ShouldBindJSON(&user)
	//JSONデータをUser構造体にバインドしてバリデーションを実行
	if err != nil {
		//バリデーションチェックを実行
		err := service.ValidationCheck(c, err)
		if err != nil {
			log.Printf("セッションidが一致しませんでした。user.ChangeId: %s, user.NowId: %s, err: %v", user.ChangeId, user.NowId, err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
	}

	//DBにIDを変更
	blog, err := db.UpdateId(user.ChangeId, user.NowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("error", err.Error())
		return
	}

	//redisでセッション破棄、新IDでセッション作成
	err = redis.UpdateSession(c, user.ChangeId, user.NowId)
	if err != nil {
		log.Println("セッション破棄と変更後IDでセッション作成に失敗しました。", err)
	}

	log.Printf("Success Change blog.UserId :blog.UserId %+v", blog.UserId)
	c.JSON(http.StatusOK, gin.H{"blog.UserId": blog.UserId})
}

// 会員情報編集(password)
func PostUpdatePw(c *gin.Context, redis redis.SessionStore) {
	user := entity.UserPwChange{}

	//リクエストをGo構造体にバインド
	err := c.ShouldBindJSON(&user)
	//JSONデータをUser構造体にバインドしてバリデーションを実行
	if err != nil {
		//バリデーションチェックを実行
		err := service.ValidationCheck(c, err)
		if err != nil {
			log.Printf("バリデーションチェックエラーが発生しました。err: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	//DBにPWを変更
	blog, err := db.UpdatePassword(user.UserId, user.NowPassword, user.ChangePassword)
	if err != nil {
		log.Printf("DBでPWを変更できませんでした。user.UserId: %s, err: %v", user.UserId, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Redisよりログイン情報セッションを一度消去
	err = redis.DeleteSession(c, user.UserId)
	if err != nil {
		log.Printf("セッションを消去できませんでした。user.UserId: %s, err: %v", user.UserId, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in decideLogout": err.Error()})
		return
	}

	//RedisよりセッションとCookieにUserIdを新しく登録
	cookieKey := os.Getenv("LOGIN_USER_ID_KEY")
	err = redis.NewSession(c, cookieKey, user.UserId)
	if err != nil {
		log.Printf("Error in NewSession PW変更画面DB上のセッションにIDの登録に失敗しました。user.UserId: %s, err: %v", user.UserId, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in redis.NewSession of postUpdatePw": err.Error()})
		return
	}

	//会員情報のPW変更に成功
	log.Printf("Success Change PW :blog.UserId %+v", blog.UserId)
	c.JSON(http.StatusOK, gin.H{"blog.UserId": blog.UserId})
}
