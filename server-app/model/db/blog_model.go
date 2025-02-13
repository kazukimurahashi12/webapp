package db

import (
	"log"
	"strconv"

	"github.com/kazukimurahashi12/webapp/model/entity"
)

// 送られてきたタイトルと内容をDBに登録
func Create(loginID, title, content string) (*entity.Blog, error) {
	blog := entity.Blog{}
	blog = entity.Blog{LoginID: loginID, Title: title, Content: content}

	result := Db.Table("BLOGS").Create(&blog)
	if result.Error != nil {
		log.Println("error", result.Error)
		// エラーが発生した場合はエラーを返す
		return nil, result.Error
	}

	return &blog, nil
}

// 送られてきたタイトルと内容をDBに更新
func Edit(id, loginID, title, content string) (*entity.Blog, error) {
	blog := entity.Blog{}
	blog = entity.Blog{LoginID: loginID, Title: title, Content: content}

	result := Db.Table("BLOGS").Where("id = ?", id).Updates(&blog)
	if result.Error != nil {
		log.Println("error", result.Error)
		// エラーが発生した場合はエラーを返す
		return nil, result.Error
	}

	return &blog, nil
}

// DBからBLOG情報を全件取得
func GetBlogOverview() ([]entity.Blog, error) {
	var blogs []entity.Blog

	//MySQLからuserIdに一致する構造体userを取得
	result := Db.Table("BLOGS").Find(&blogs)
	if result.Error != nil {
		log.Println("error", result.Error)
		// エラーが発生した場合はエラーを返す
		return nil, result.Error
	}

	return blogs, nil
}

// DBからIDによる特定のBLOG情報を取得
func GetBlogViewInfoById(id string) (*entity.Blog, error) {
	blog := &entity.Blog{}

	// MySQLからIDに一致する構造体blogを取得
	result := Db.Table("BLOGS").Where("id = ?", id).First(blog)
	if result.Error != nil {
		log.Println("error", result.Error)
		// エラーが発生した場合はエラーを返す
		return nil, result.Error
	}
	return blog, nil
}

// DBからIDによる特定のBLOG情報を消去
func DeleteBlogInfoById(id string) (*entity.Blog, error) {
	// 文字列をuintに変換
	uintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}
	blog := &entity.Blog{ID: uint(uintId)} //リクセストIDをプロパティに格納

	//消去予定のBLOG情報を取得
	result := Db.Table("BLOGS").Where("id = ?", id).First(blog)
	if result.Error != nil {
		log.Println("Error retrieving blog:", result.Error)
		return nil, result.Error
	}
	//消去対象のBLOG情報
	deletedBlog := *blog

	// MySQLからIDに一致する構造体blogを取得
	deleteResult := Db.Table("BLOGS").Delete(blog)

	if deleteResult.Error != nil {
		log.Println("error", deleteResult.Error)
		// エラーが発生した場合はエラーを返す
		return nil, deleteResult.Error
	}

	//成功消去の場合、消去されたBLOG情報をログ出力
	log.Println("Deleted blog:ID, Title", deletedBlog.ID, deletedBlog.Title)
	return blog, nil
}
