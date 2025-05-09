package blog

import (
	"errors"
)

// DB保存用のBlog を生成するファクトリ関数
// ID（uint）CreatedAt/UpdatedAt自動生成
func NewBlog(authorID uint, title, content string) (*Blog, error) {
	if title == "" || len(title) > 50 {
		return nil, errors.New("タイトルの長さが不正です")
	}
	if content == "" || len(content) > 8000 {
		return nil, errors.New("本文の長さが不正です")
	}

	return &Blog{
		AuthorID: authorID,
		Title:    title,
		Content:  content,
		// CreatedAt, UpdatedAt GORM自動で設定
	}, nil
}
