package blog

import (
	"strconv"

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

func (b *blogUseCase) NewCreateBlog(blog *domainBlog.BlogPost) (*domainBlog.BlogPost, error) {
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

func (b *blogUseCase) UpdateBlog(blogPost *domainBlog.BlogPost) (*domainBlog.BlogPost, error) {
	// stringからuintに変換
	id, err := strconv.ParseUint(blogPost.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	blog := &domainBlog.Blog{
		ID:      uint(id),
		LoginID: blogPost.LoginID,
		Title:   blogPost.Title,
		Content: blogPost.Content,
	}

	updateErr := b.blogRepo.Update(blog)
	if updateErr != nil {
		return nil, updateErr
	}

	return &domainBlog.BlogPost{
		ID:      strconv.FormatUint(uint64(blog.ID), 10),
		LoginID: blog.LoginID,
		Title:   blog.Title,
		Content: blog.Content,
	}, nil
}

func (b *blogUseCase) NewCreateUser(user *domainUser.FormUser) (*domainUser.User, error) {
	// FormUserをUserに変換
	newUser := &domainUser.User{
		UserID:   user.UserID,
		Password: user.Password,
	}

	// ユーザー登録処理
	err := b.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (b *blogUseCase) GetBlogByID(id string) (*domainBlog.Blog, error) {
	blog, err := b.blogRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return blog, nil
}
