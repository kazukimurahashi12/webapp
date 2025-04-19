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

func (b *blogUseCase) GetBlogsByUserID(userID string) ([]domainBlog.Blog, error) {
	return b.blogRepo.FindByUserID(userID)
}

func (b *blogUseCase) GetUserByID(userID string) (*domainUser.User, error) {
	return b.userRepo.FindByID(userID)
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

func (b *blogUseCase) NewCreateUser(user *domainUser.User) (*domainUser.User, error) {
	// ユーザー登録処理
	err := b.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b *blogUseCase) GetBlogByID(id string) (*domainBlog.Blog, error) {
	blog, err := b.blogRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return blog, nil
}
