package blog

import (
	"errors"
	"time"
)

// DB保存用のBlog を生成するファクトリ関数
// ID（uint）CreatedAt/UpdatedAt自動生成
func NewBlog(post *BlogPost) (*Blog, error) {
	if len(post.LoginID) < 2 || len(post.LoginID) > 10 {
		return nil, errors.New("LoginIDの長さが不正です")
	}
	if len(post.Title) < 1 || len(post.Title) > 50 {
		return nil, errors.New("タイトルの長さが不正です")
	}
	if len(post.Content) < 1 || len(post.Content) > 8000 {
		return nil, errors.New("本文の長さが不正です")
	}

	// 現在時刻を取得
	now := time.Now()

	//GORMが自動でインクリメントID
	return &Blog{
		LoginID:   post.LoginID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
