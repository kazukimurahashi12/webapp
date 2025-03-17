package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/infrastructure/db"
	"github.com/kazukimurahashi12/webapp/usecase/validator"
)

// 新規会員登録(id,password)
func PostRegist(c *gin.Context) {
	//構造体をインスタンス化
	registUser := domain.FormUser{}
	//リクエストをGo構造体にバインド
	err := c.ShouldBindJSON(&registUser)
	//JSONデータをUser構造体にバインドしてバリデーションを実行
	if err != nil {
		//バリデーションチェックを実行
		err := validator.ValidationCheck(c, err)
		if err != nil {
			log.Printf("リクエストJSON形式で構造体にバインドを失敗しました。registUser.UserId: %s, err: %v", registUser.UserId, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	//DBクラインアントのインスタンスを作成
	dbClient := &db.DB{}
	//userRepositoryのインスタンスを作成
	userRepository := db.NewUserRepository(dbClient)
	//DBに会員情報登録処理
	user := &domain.User{
		UserId:   registUser.UserId,
		Password: registUser.Password,
	}
	err = userRepository.Create(user)
	if err != nil {
		log.Printf("DBに会員情報の登録に失敗しました。registUser.UserId: %s, err: %v", registUser.UserId, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error in db.Signup of postRegist": err.Error()})
		return
	}
	//DBに会員情報登録に成功
	log.Printf("Success user in RegisterView from DB :user %+v", user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}
