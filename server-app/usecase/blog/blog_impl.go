package blog

import (
	"strconv"

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

func (b *blogUseCase) DeleteBlog(id string) error {
	return b.blogRepo.Delete(id)
}

func (b *blogUseCase) UpdateBlog(blogPost *domain.BlogPost) (*domain.BlogPost, error) {
	// stringからuintに変換
	id, err := strconv.ParseUint(blogPost.Id, 10, 64)
	if err != nil {
		return nil, err
	}

	blog := &domain.Blog{
		Id:      uint(id),
		LoginId: blogPost.LoginId,
		Title:   blogPost.Title,
		Content: blogPost.Content,
	}

	updateErr := b.blogRepo.Update(blog)
	if updateErr != nil {
		return nil, updateErr
	}

	return &domain.BlogPost{
		Id:      strconv.FormatUint(uint64(blog.Id), 10),
		LoginId: blog.LoginId,
		Title:   blog.Title,
		Content: blog.Content,
	}, nil
}

func (b *blogUseCase) NewCreateUser(user *domain.FormUser) (*domain.User, error) {
	// FormUserをUserに変換
	newUser := &domain.User{
		UserId:   user.UserId,
		Password: user.Password,
	}

	// ユーザー登録処理
	err := b.userRepo.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (b *blogUseCase) GetBlogByID(id string) (*domain.Blog, error) {
	blog, err := b.blogRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return blog, nil
}
