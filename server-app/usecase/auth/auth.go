package auth

import domainUser "github.com/kazukimurahashi12/webapp/domain/user"

type UseCase interface {
	Authenticate(loginID, password string) (*domainUser.User, error)
	GetUserByID(loginID string) (*domainUser.User, error)
}
