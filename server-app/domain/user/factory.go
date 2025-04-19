package user

import (
	"errors"
	"time"
)

// FormUserからUserに変換するファクトリ関数
func NewUser(userID, password string) (*User, error) {
	if len(userID) < 2 || len(userID) > 10 {
		return nil, errors.New("ユーザーIDの長さが不正です")
	}
	if len(password) < 4 || len(password) > 20 {
		return nil, errors.New("パスワードの長さが不正です")
	}

	now := time.Now()

	return &User{
		UserID:    userID,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
