package mapper

import (
	"github.com/kazukimurahashi12/webapp/domain/blog"
	"github.com/kazukimurahashi12/webapp/interface/dto"
)

func ToBlogCreatedResponse(b *blog.Blog) *dto.BlogCreatedResponse {
	return &dto.BlogCreatedResponse{
		ID:    b.ID,
		Title: b.Title,
	}
}
