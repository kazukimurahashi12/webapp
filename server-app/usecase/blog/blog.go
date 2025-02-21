package blog

import "github.com/kazukimurahashi12/webapp/domain"

type UseCase interface {
	GetBlogsByUserID(userID string) ([]domain.Blog, error)
	GetUserByID(userID string) (*domain.User, error)
}
