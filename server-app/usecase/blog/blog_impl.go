package blog

import (
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
)

type blogUseCase struct {
	blogRepo domainBlog.BlogRepository
	userRepo domainUser.UserRepository
}

func NewBlogUseCase(blogRepo domainBlog.BlogRepository, userRepo domainUser.UserRepository) UseCase {
	return &blogUseCase{
		blogRepo: blogRepo,
		userRepo: userRepo,
	}
}

func (b *blogUseCase) NewCreateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error) {
	err := b.blogRepo.Create(blog)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

func (b *blogUseCase) FindBlogsByUserID(userID string) ([]domainBlog.Blog, error) {
	return b.blogRepo.FindBlogsByUserID(userID)
}

func (b *blogUseCase) FindBlogByID(id uint) (*domainBlog.Blog, error) {
	return b.blogRepo.FindBlogByID(id)
}

func (b *blogUseCase) DeleteBlog(id string) error {
	return b.blogRepo.Delete(id)
}

func (b *blogUseCase) UpdateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error) {

	updateErr := b.blogRepo.Update(blog)
	if updateErr != nil {
		return nil, updateErr
	}

	return blog, nil
}
