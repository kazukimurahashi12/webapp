package auth

import (
	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
)

type authUseCase struct {
	userRepo domainUser.UserRepository
}

func NewAuthUseCase(userRepo domainUser.UserRepository) UseCase {
	return &authUseCase{
		userRepo: userRepo,
	}
}

// ユーザーIDとパスワードを元に認証
func (a *authUseCase) Authenticate(userID, password string) (*domainUser.User, error) {
	return a.userRepo.FindByID(userID)
}

// ユーザーIDを元にユーザー情報を取得
func (a *authUseCase) GetUserByID(userID string) (*domainUser.User, error) {
	return a.userRepo.FindByID(userID)
}
