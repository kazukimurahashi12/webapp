package user

import "github.com/kazukimurahashi12/webapp/domain"

type UseCase interface {
	UpdateUserID(oldID, newID string) (*domain.User, error)
	UpdateUserPassword(userID, currentPassword, newPassword string) (*domain.User, error)
}
