package blog

import "github.com/kazukimurahashi12/webapp/domain"

type UseCase interface {
	NewCreateBlog(blog *domain.BlogPost) (*domain.BlogPost, error)
	GetBlogsByUserID(userID string) ([]domain.Blog, error)
	GetUserByID(userID string) (*domain.User, error)
}
