package auth

import (
	"strconv"

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
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, err
	}
	return a.userRepo.FindUserByUserID(uint(userIDUint))
}

// ユーザーIDを元にユーザー情報を取得
func (a *authUseCase) GetUserByID(userID string) (*domainUser.User, error) {
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, err
	}
	return a.userRepo.FindUserByUserID(uint(userIDUint))
}
