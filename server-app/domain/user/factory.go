package user

import (
	"errors"
	"time"
)

// FormUserからUserに変換するファクトリ関数
func NewUser(form *FormUser) (*User, error) {
	if len(form.UserID) < 2 || len(form.UserID) > 10 {
		return nil, errors.New("ユーザーIDの長さが不正です")
	}
	if len(form.Password) < 4 || len(form.Password) > 20 {
		return nil, errors.New("パスワードの長さが不正です")
	}

	now := time.Now()

	return &User{
		UserID:    form.UserID,
		Password:  form.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
