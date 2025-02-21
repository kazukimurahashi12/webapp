package user

import (
	"github.com/kazukimurahashi12/webapp/crypto"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/repository"
)

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) UpdateUserID(oldID, newID string) (*domain.User, error) {
	return u.userRepo.UpdateID(oldID, newID)
}

func (u *userUseCase) UpdateUserPassword(userID, currentPassword, newPassword string) (*domain.User, error) {
	// 現在のパスワードを検証
	user, err := u.userRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if err := crypto.CompareHashAndPassword(user.Password, currentPassword); err != nil {
		return nil, err
	}

	return u.userRepo.UpdatePassword(userID, newPassword)
}
