package user

import (
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/infrastructure/crypto"
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

	crypto := crypto.NewBcryptCrypto()
	if err := crypto.CompareHashAndPassword(user.Password, currentPassword); err != nil {
		return nil, err
	}

	return u.userRepo.UpdatePassword(userID, newPassword)
}

func (u *userUseCase) CreateUser(user *domain.FormUser) (*domain.User, error) {
	// パスワードをハッシュ化
	crypto := crypto.NewBcryptCrypto()
	hashedPassword, err := crypto.Encrypt(user.Password)
	if err != nil {
		return nil, err
	}

	// 新しいユーザーを作成
	newUser := &domain.User{
		UserId:   user.UserId,
		Password: hashedPassword,
	}

	// ユーザーを登録
	err = u.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
