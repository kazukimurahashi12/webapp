package blog

import (
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/repository"
)

type blogUseCase struct {
	blogRepo repository.BlogRepository
	userRepo repository.UserRepository
}

func NewBlogUseCase(blogRepo repository.BlogRepository, userRepo repository.UserRepository) UseCase {
	return &blogUseCase{
		blogRepo: blogRepo,
		userRepo: userRepo,
	}
}

func (b *blogUseCase) NewCreateBlog(blog *domain.BlogPost) (*domain.BlogPost, error) {
	err := b.blogRepo.Create(blog)
	if err != nil {
		return nil, err
	}
	return blog, nil
}
func (b *blogUseCase) GetBlogsByUserID(userID string) ([]domain.Blog, error) {
	return b.blogRepo.FindByUserID(userID)
}

func (b *blogUseCase) GetUserByID(userID string) (*domain.User, error) {
	return b.userRepo.FindByID(userID)
}
