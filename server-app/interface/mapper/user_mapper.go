package mapper

import (
	"github.com/kazukimurahashi12/webapp/domain/user"
	"github.com/kazukimurahashi12/webapp/interface/dto"
)

func ToUserCreatedResponse(u *user.User) *dto.UserCreatedResponse {
	return &dto.UserCreatedResponse{
		ID:        u.ID,
		UserID:    u.Username,
		CreatedAt: u.CreatedAt,
	}
}

func ToUserIDResponse(u *user.User) *dto.UserIDResponse {
	return &dto.UserIDResponse{
		UserID: u.Username,
	}
}
