package blog

import "github.com/kazukimurahashi12/webapp/domain"

type UseCase interface {
	NewCreateBlog(blog *domain.BlogPost) (*domain.BlogPost, error)
	GetBlogsByUserID(userID string) ([]domain.Blog, error)
	GetUserByID(userID string) (*domain.User, error)
	DeleteBlog(id string) error
	UpdateBlog(blog *domain.BlogPost) (*domain.BlogPost, error)
	GetBlogByID(id string) (*domain.Blog, error)
	NewCreateUser(user *domain.FormUser) (*domain.User, error)
}
