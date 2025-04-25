package user

import domainUser "github.com/kazukimurahashi12/webapp/domain/user"

type UseCase interface {
	FindUserByUserID(userID string) (*domainUser.User, error)
	FindUserByID(id uint) (*domainUser.User, error)
	UpdateUserID(oldID, newID string) (*domainUser.User, error)
	UpdateUserPassword(userID, currentPassword, newPassword string) (*domainUser.User, error)
	CreateUser(userID, password string) (*domainUser.User, error)
}
