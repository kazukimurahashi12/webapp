package blog

import "errors"

// ドメインエラーの定義
var (
	ErrBlogNotFound        = errors.New("blog not found")
	ErrBlogAlreadyExists   = errors.New("blog with the same title already exists")
	ErrBlogInvalidData     = errors.New("blog data is invalid")
	ErrBlogUnauthorized    = errors.New("unauthorized access to this blog")
	ErrBlogContentEmpty    = errors.New("blog content cannot be empty")
	ErrBlogTitleTooLong    = errors.New("blog title exceeds maximum length")
	ErrBlogTitleEmpty      = errors.New("blog title cannot be empty")
	ErrBlogVersionConflict = errors.New("blog has been modified by another user")
	ErrBlogDeleted         = errors.New("blog has been deleted")
	ErrBlogPublishFailed   = errors.New("failed to publish blog")
)
