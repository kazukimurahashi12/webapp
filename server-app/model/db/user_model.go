package db

import (
	"errors"
	"log"

	"github.com/kazukimurahashi12/webapp/crypto"
	"github.com/kazukimurahashi12/webapp/model/entity"
)

func init() {
	//MySQLのストレージエンジンInnoDB,テーブル自動生成
	Db.Table("BLOGS").Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(entity.Blog{})
	Db.Table("USERS").Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(entity.User{})
}

// 送られてきたUserIdと一致するUserを取得
// 取得したUserの暗号化済みPasswordと送られてきたPasswordをCompareHashAndPassword関数でチェック
func CheckUser(userId, password string) (*entity.User, error) {
	user := entity.User{}
	//MySQLからuserIdに一致する構造体userを取得
	Db.Table("USERS").Where("user_id = ?", userId).First(&user)

	if user.ID == 0 {
		err := errors.New("ユーザー名が一致しません。")
		log.Println(err)
		return nil, err
	}

	//ハッシュ化したpasswordを比較
	compareErr := crypto.CompareHashAndPassword(user.Password, password)

	if compareErr != nil {
		err := errors.New("パスワードが一致しません。:" + compareErr.Error())
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

// 送られてきたUserId,Passwordと一致するUserが既に登録されているか確認
// PasswordEncrypt関数で送られてきたPasswordを暗号化
// gormのCreate関数で新規会員登録
func Signup(userId, password string) (*entity.User, error) {
	user := entity.User{}

	Db.Table("USERS").Where("user_id = ?", userId).First(&user)

	if user.ID != 0 {
		err := errors.New("同一名のUserIdが既に登録されています。")
		log.Println(err)
		return nil, err
	}

	//ハッシュ化したpasswordを作成
	encryptPw, err := crypto.PasswordEncrypt(password)

	if err != nil {
		log.Println("パスワード暗号化中にエラーが発生しました。：", err)
		return nil, err
	}

	user = entity.User{UserId: userId, Password: encryptPw}
	Db.Table("USERS").Create(&user)
	return &user, nil
}

// gormのUpdate関数でID情報編集
func UpdateId(changeId string, nowId string) (*entity.User, error) {
	user := entity.User{}

	if err := Db.Table("USERS").Where("user_id = ?", changeId).First(&user).Error; err == nil {
		log.Println("UserIdが重複するユーザーが存在しています。")
		log.Println("Error duplicate id from DB", user)
		return nil, errors.New("UserIdが重複するユーザーが存在しています。")
	} else {
		log.Println("重複するユーザーはDBに存在しません。", err)
		log.Println("Request id from client", changeId)
	}

	// 既存のUSER情報を取得
	newUser := entity.User{}
	if err := Db.Table("USERS").Where("user_id = ?", nowId).First(&newUser).Error; err != nil {
		log.Printf("USER情報の取得に失敗しました。nowId: %s, err: %v", nowId, err.Error())
		return nil, err
	}

	//指定されたフィールドのみを更新
	if err := Db.Table("USERS").Where("user_id = ?", nowId).Updates(&newUser).Error; err != nil {
		err := errors.New("IDの更新に失敗しました。")
		log.Println(err)
		return nil, err
	}

	//成功消去の場合、消去されたBLOG情報をログ出力
	log.Println("IDの変更に成功しました。", newUser.UserId)
	return &newUser, nil
}

// gormのUpdate関数でPW情報編集
func UpdatePassword(userId, nowPassword string, newPassword string) (*entity.User, error) {
	user := entity.User{}
	//MySQLからuserIdに一致する構造体userを取得
	Db.Table("USERS").Where("user_id = ?", userId).First(&user)
	if user.ID == 0 {
		err := errors.New("ユーザー名が一致しません。")
		log.Println(err)
		return nil, err
	}

	//DB登録されているPWとリクエストされたPWを比較
	compareErr := crypto.CompareHashAndPassword(user.Password, nowPassword)
	if compareErr != nil {
		err := errors.New("パスワードが一致しません。:" + compareErr.Error())
		log.Println(err)
		return nil, err
	}

	//NewPaswweordの設定
	//ハッシュ化したpasswordを作成
	encryptPw, err := crypto.PasswordEncrypt(newPassword)
	if err != nil {
		log.Println("パスワード暗号化中にエラーが発生しました。：", err)
		return nil, err
	}

	// userPw := entity.User{UserId: userId, Password: encryptPw}
	user.Password = encryptPw

	//指定されたフィールドのみを更新
	if err := Db.Table("USERS").Where("user_id = ?", userId).Updates(&user).Error; err != nil {
		log.Printf("USER情報の取得に失敗しました。userId: %s, errChangePw: %v", userId, err.Error())
		return nil, err
	}

	return &user, nil
}

// gormのUpdate関数でuserIdに一致する構造体userを取得
func GetOneUser(UserId string) (entity.User, error) {
	user := entity.User{}
	//MySQLからuserIdに一致する構造体userを取得
	result := Db.Table("USERS").Where("user_id = ?", UserId).First(&user)
	if result.Error != nil {
		log.Println("error", result.Error)
		// エラーが発生した場合はエラーを返す
		return user, result.Error
	}
	return user, nil
}
