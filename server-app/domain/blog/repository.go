package blog

// ブログRepositoryインターフェース
type BlogRepository interface {
	Create(blog *Blog) error
	FindBlogByID(id uint) (*Blog, error)
	FindBlogsByUserID(userID string) ([]Blog, error)
	Update(blog *Blog) error
	Delete(id string) error
}
