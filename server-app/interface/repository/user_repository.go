package repository

import (
	"github.com/kazukimurahashi12/webapp/domain"
)

type UserRepository interface {
	FindByID(id string) (*domain.User, error)
	FindByUserID(userID string) (*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	UpdateID(oldID, newID string) (*domain.User, error)
	UpdatePassword(userID, newPassword string) (*domain.User, error)
}
