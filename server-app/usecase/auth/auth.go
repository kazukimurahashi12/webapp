package auth

import domainUser "github.com/kazukimurahashi12/webapp/domain/user"

type UseCase interface {
	Authenticate(userID, password string) (*domainUser.User, error)
	GetUserByID(userID string) (*domainUser.User, error)
}
