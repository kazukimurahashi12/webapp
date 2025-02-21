package auth

import "github.com/kazukimurahashi12/webapp/domain"

type UseCase interface {
	Authenticate(userID, password string) (*domain.User, error)
	GetUserByID(userID string) (*domain.User, error)
}
