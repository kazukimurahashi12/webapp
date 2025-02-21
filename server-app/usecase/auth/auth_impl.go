package auth

import (
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/repository"
)

type authUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) UseCase {
	return &authUseCase{
		userRepo: userRepo,
	}
}

func (a *authUseCase) Authenticate(userID, password string) (*domain.User, error) {
	return a.userRepo.FindByID(userID)
}

func (a *authUseCase) GetUserByID(userID string) (*domain.User, error) {
	return a.userRepo.FindByID(userID)
}
