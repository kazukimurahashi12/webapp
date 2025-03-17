package repository

import "github.com/kazukimurahashi12/webapp/domain"

type BlogRepository interface {
	Create(blog *domain.BlogPost) error
	FindByID(id string) (*domain.Blog, error)
	FindByUserID(userID string) ([]domain.Blog, error)
	Update(blog *domain.Blog) error
	Delete(id string) error
}
