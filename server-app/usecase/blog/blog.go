package blog

import (
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
)

type UseCase interface {
	NewCreateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error)
	FindBlogsByUserID(userID string) ([]domainBlog.Blog, error)
	FindBlogByID(id uint) (*domainBlog.Blog, error)
	DeleteBlog(id string) error
	UpdateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error)
}
