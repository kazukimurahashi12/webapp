package blog

import (
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
)

type UseCase interface {
	NewCreateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error)
	FindBlogsByAuthorID(authorID uint) ([]domainBlog.Blog, error)
	FindBlogByID(id uint) (*domainBlog.Blog, error)
	FindBlogByAuthorID(authorID uint) (*domainBlog.Blog, error)
	DeleteBlog(id uint) error
	UpdateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error)
}
