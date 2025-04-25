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

func ToBlogsResponse(blogs []blog.Blog) []*dto.BlogCreatedResponse {
	responses := make([]*dto.BlogCreatedResponse, len(blogs))

	for i, b := range blogs {
		responses[i] = &dto.BlogCreatedResponse{
			ID:    b.ID,
			Title: b.Title,
		}
	}

	return responses
}
