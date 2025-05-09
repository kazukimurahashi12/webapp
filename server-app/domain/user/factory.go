package user

import (
	"errors"
	"time"
)

// FormUserからUserに変換するファクトリ関数
func NewUser(username, password string) (*User, error) {
	if len(username) < 2 || len(username) > 10 {
		return nil, errors.New("ユーザー名の長さが不正です")
	}
	if len(password) < 4 || len(password) > 20 {
		return nil, errors.New("パスワードの長さが不正です")
	}

	now := time.Now()

	return &User{
		Username:  username,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
