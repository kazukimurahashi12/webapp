package blog

import (
	"errors"
	"time"
)

// DB保存用のBlog を生成するファクトリ関数
// ID（uint）CreatedAt/UpdatedAt自動生成
func NewBlog(loginID, title, content string) (*Blog, error) {
	// バリデーションチェック
	if len(loginID) < 2 || len(loginID) > 10 {
		return nil, errors.New("LoginIDの長さが不正です")
	}
	if len(title) < 1 || len(title) > 50 {
		return nil, errors.New("タイトルの長さが不正です")
	}
	if len(content) < 1 || len(content) > 8000 {
		return nil, errors.New("本文の長さが不正です")
	}

	// 現在時刻を取得
	now := time.Now()

	// GORMが自動でインクリメントID
	return &Blog{
		LoginID:   loginID,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
