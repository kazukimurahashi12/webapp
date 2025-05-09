package blog

// ブログRepositoryインターフェース
type BlogRepository interface {
	Create(blog *Blog) error
	FindBlogByID(id uint) (*Blog, error)
	FindBlogsByAuthorID(authorID uint) ([]Blog, error)
	FindBlogByAuthorID(authorID uint) (*Blog, error)
	Update(blog *Blog) error
	Delete(id uint) error
}
