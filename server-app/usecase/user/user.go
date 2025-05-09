package user

import domainUser "github.com/kazukimurahashi12/webapp/domain/user"

type UseCase interface {
	FindUserByUserID(userID uint) (*domainUser.User, error)
	FindUserByID(id uint) (*domainUser.User, error)
	UpdateUserID(oldID, newID uint) (*domainUser.User, error)
	UpdateUserPassword(userID uint, currentPassword, newPassword string) (*domainUser.User, error)
	CreateUser(username, password string) (*domainUser.User, error)
}
