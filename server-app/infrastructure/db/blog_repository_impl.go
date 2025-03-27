package db

import (
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type blogRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewBlogRepository(manager *DBManager) repository.BlogRepository {
	return &blogRepository{
		db:     manager.db,
		logger: manager.logger,
	}
}

// ブログを作成
func (r *blogRepository) Create(blog *domain.BlogPost) error {
	newBlog := domain.Blog{
		Title:   blog.Title,
		Content: blog.Content,
		LoginId: blog.LoginId,
	}
	return r.db.Table("BLOGS").Create(&newBlog).Error
}

// ブログを取得
func (r *blogRepository) FindByID(id string) (*domain.Blog, error) {
	blog := domain.Blog{}
	if err := r.db.Table("BLOGS").Where("id = ?", id).First(&blog).Error; err != nil {
		return nil, err
	}
	return &domain.Blog{
		Id:      blog.Id,
		Title:   blog.Title,
		Content: blog.Content,
		LoginId: blog.LoginId,
	}, nil
}

// ユーザーIDに紐づくブログを取得
func (r *blogRepository) FindByUserID(userID string) ([]domain.Blog, error) {
	var blogs []domain.Blog
	if err := r.db.Table("BLOGS").Where("user_id = ?", userID).Find(&blogs).Error; err != nil {
		return nil, err
	}

	var domainBlogs []domain.Blog
	for _, blog := range blogs {
		domainBlogs = append(domainBlogs, domain.Blog{
			Id:      blog.Id,
			Title:   blog.Title,
			Content: blog.Content,
			LoginId: blog.LoginId,
		})
	}
	return domainBlogs, nil
}

// ブログを更新
func (r *blogRepository) Update(blog *domain.Blog) error {
	existingBlog := domain.Blog{}
	if err := r.db.Table("BLOGS").Where("id = ?", blog.Id).First(&existingBlog).Error; err != nil {
		return err
	}

	existingBlog.Title = blog.Title
	existingBlog.Content = blog.Content
	return r.db.Table("BLOGS").Save(&existingBlog).Error
}

// ブログを削除
func (r *blogRepository) Delete(id string) error {
	return r.db.Table("BLOGS").Where("id = ?", id).Delete(&domain.Blog{}).Error
}
