package user

import "errors"

// ドメインエラーの定義
var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user with this ID already exists")
	ErrInvalidUserID        = errors.New("invalid user ID format")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrUnauthorized         = errors.New("unauthorized access")
	ErrPasswordTooWeak      = errors.New("password does not meet strength requirements")
	ErrUserLocked           = errors.New("user account is locked")
	ErrUserDisabled         = errors.New("user account is disabled")
	ErrAuthenticationFailed = errors.New("authentication failed")
)
