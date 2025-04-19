package mapper

import (
	"github.com/kazukimurahashi12/webapp/domain/user"
	"github.com/kazukimurahashi12/webapp/interface/dto"
)

func ToUserCreatedResponse(u *user.User) *dto.UserCreatedResponse {
	return &dto.UserCreatedResponse{
		ID:        u.ID,
		UserID:    u.UserID,
		CreatedAt: u.CreatedAt,
	}
}
