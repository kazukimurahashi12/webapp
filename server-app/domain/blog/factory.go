package blog

import (
	"errors"
)

// DB保存用のBlog を生成するファクトリ関数
// ID（uint）CreatedAt/UpdatedAt自動生成
// func NewBlog(userID, title, content string) (*Blog, error) {
// 	// バリデーションチェック
// 	if len(userID) < 2 || len(userID) > 10 {
// 		return nil, errors.New("LoginIDの長さが不正です")
// 	}
// 	if len(title) < 1 || len(title) > 50 {
// 		return nil, errors.New("タイトルの長さが不正です")
// 	}
// 	if len(content) < 1 || len(content) > 8000 {
// 		return nil, errors.New("本文の長さが不正です")
// 	}

// 	// 現在時刻を取得
// 	now := time.Now()

// 	// GORMが自動でインクリメントID
// 	return &Blog{
// 		UserID:    userID,
// 		Title:     title,
// 		Content:   content,
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}, nil
// }

func NewBlog(userID uint, title, content string) (*Blog, error) {
	if title == "" || len(title) > 50 {
		return nil, errors.New("タイトルの長さが不正です")
	}
	if content == "" || len(content) > 8000 {
		return nil, errors.New("本文の長さが不正です")
	}

	return &Blog{
		UserID:  userID,
		Title:   title,
		Content: content,
		// CreatedAt, UpdatedAt GORM自動で設定
	}, nil
}
