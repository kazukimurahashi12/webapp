package blog

import (
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
)

type UseCase interface {
	NewCreateBlog(blog *domainBlog.BlogPost) (*domainBlog.BlogPost, error)
	GetBlogsByUserID(userID string) ([]domainBlog.Blog, error)
	GetUserByID(userID string) (*domainUser.User, error)
	DeleteBlog(id string) error
	UpdateBlog(blog *domainBlog.BlogPost) (*domainBlog.BlogPost, error)
	GetBlogByID(id string) (*domainBlog.Blog, error)
	NewCreateUser(user *domainUser.FormUser) (*domainUser.User, error)
}
