package blog

import (
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
)

type UseCase interface {
	NewCreateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error)
	GetBlogsByUserID(userID string) ([]domainBlog.Blog, error)
	GetUserByID(userID string) (*domainUser.User, error)
	DeleteBlog(id string) error
	UpdateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error)
	GetBlogByID(id string) (*domainBlog.Blog, error)
	NewCreateUser(user *domainUser.User) (*domainUser.User, error)
}
