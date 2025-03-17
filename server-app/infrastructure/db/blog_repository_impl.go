package db

import (
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/repository"
)

type blogRepository struct {
	db *DB
}

func NewBlogRepository(db *DB) repository.BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) Create(blog *domain.BlogPost) error {
	newBlog := domain.Blog{
		Title:   blog.Title,
		Content: blog.Content,
		LoginID: blog.LoginID,
	}
	return r.db.Table("BLOGS").Create(&newBlog).Error
}

func (r *blogRepository) FindByID(id string) (*domain.Blog, error) {
	blog := domain.Blog{}
	if err := r.db.Table("BLOGS").Where("id = ?", id).First(&blog).Error; err != nil {
		return nil, err
	}
	return &domain.Blog{
		ID:      blog.ID,
		Title:   blog.Title,
		Content: blog.Content,
		LoginID: blog.LoginID,
	}, nil
}

func (r *blogRepository) FindByUserID(userID string) ([]domain.Blog, error) {
	var blogs []domain.Blog
	if err := r.db.Table("BLOGS").Where("user_id = ?", userID).Find(&blogs).Error; err != nil {
		return nil, err
	}

	var domainBlogs []domain.Blog
	for _, blog := range blogs {
		domainBlogs = append(domainBlogs, domain.Blog{
			ID:      blog.ID,
			Title:   blog.Title,
			Content: blog.Content,
			LoginID: blog.LoginID,
		})
	}
	return domainBlogs, nil
}

func (r *blogRepository) Update(blog *domain.Blog) error {
	existingBlog := domain.Blog{}
	if err := r.db.Table("BLOGS").Where("id = ?", blog.ID).First(&existingBlog).Error; err != nil {
		return err
	}

	existingBlog.Title = blog.Title
	existingBlog.Content = blog.Content
	return r.db.Table("BLOGS").Save(&existingBlog).Error
}

func (r *blogRepository) Delete(id string) error {
	return r.db.Table("BLOGS").Where("id = ?", id).Delete(&domain.Blog{}).Error
}
