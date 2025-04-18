package blog

// ブログRepositoryインターフェース
type BlogRepository interface {
	Create(blog *BlogPost) error
	FindByID(id string) (*Blog, error)
	FindByUserID(userID string) ([]Blog, error)
	Update(blog *Blog) error
	Delete(id string) error
}
