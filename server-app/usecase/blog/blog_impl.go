package blog

import (
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
)

type blogUseCase struct {
	blogRepo domainBlog.BlogRepository
}

func NewBlogUseCase(blogRepo domainBlog.BlogRepository) UseCase {
	return &blogUseCase{
		blogRepo: blogRepo,
	}
}

func (b *blogUseCase) NewCreateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error) {
	err := b.blogRepo.Create(blog)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

func (b *blogUseCase) FindBlogsByAuthorID(authorID uint) ([]domainBlog.Blog, error) {
	return b.blogRepo.FindBlogsByAuthorID(authorID)
}

func (b *blogUseCase) FindBlogByID(id uint) (*domainBlog.Blog, error) {
	return b.blogRepo.FindBlogByID(id)
}

func (b *blogUseCase) FindBlogByAuthorID(authorID uint) (*domainBlog.Blog, error) {
	return b.blogRepo.FindBlogByAuthorID(authorID)
}

func (b *blogUseCase) DeleteBlog(id uint) error {
	return b.blogRepo.Delete(id)
}

func (b *blogUseCase) UpdateBlog(blog *domainBlog.Blog) (*domainBlog.Blog, error) {

	updateErr := b.blogRepo.Update(blog)
	if updateErr != nil {
		return nil, updateErr
	}

	return blog, nil
}
