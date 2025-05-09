package user

import (
	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
	"github.com/kazukimurahashi12/webapp/infrastructure/crypto"
)

type userUseCase struct {
	userRepo domainUser.UserRepository
}

func NewUserUseCase(userRepo domainUser.UserRepository) UseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) FindUserByUserID(userID uint) (*domainUser.User, error) {
	return u.userRepo.FindUserByUserID(userID)
}

func (u *userUseCase) FindUserByID(id uint) (*domainUser.User, error) {
	return u.userRepo.FindUserByID(id)
}

func (u *userUseCase) UpdateUserID(oldID, newID uint) (*domainUser.User, error) {
	return u.userRepo.UpdateID(oldID, newID)
}

func (u *userUseCase) UpdateUserPassword(userID uint, currentPassword, newPassword string) (*domainUser.User, error) {
	// 現在のパスワードを検証
	user, err := u.userRepo.FindUserByUserID(userID)
	if err != nil {
		return nil, err
	}

	crypto := crypto.NewBcryptCrypto()
	if err := crypto.CompareHashAndPassword(user.Password, currentPassword); err != nil {
		return nil, err
	}

	return u.userRepo.UpdatePassword(userID, newPassword)
}

func (u *userUseCase) CreateUser(username, password string) (*domainUser.User, error) {
	// パスワードをハッシュ化
	crypto := crypto.NewBcryptCrypto()
	hashedPassword, err := crypto.Encrypt(password)
	if err != nil {
		return nil, err
	}

	// 新しいユーザーを作成
	newUser := &domainUser.User{
		Username: username,
		Password: hashedPassword,
	}

	// ユーザーを登録
	err = u.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
