package user

import domainUser "github.com/kazukimurahashi12/webapp/domain/user"

type UseCase interface {
	UpdateUserID(oldID, newID string) (*domainUser.User, error)
	UpdateUserPassword(userID, currentPassword, newPassword string) (*domainUser.User, error)
	CreateUser(user *domainUser.FormUser) (*domainUser.User, error)
}
