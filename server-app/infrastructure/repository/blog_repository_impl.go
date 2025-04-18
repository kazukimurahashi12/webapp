package repository

import (
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
	"github.com/kazukimurahashi12/webapp/infrastructure/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type blogRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewBlogRepository(manager *db.DBManager) domainBlog.BlogRepository {
	return &blogRepository{
		db:     manager.DB,
		logger: manager.Logger,
	}
}

// ブログを作成
func (r *blogRepository) Create(blog *domainBlog.BlogPost) error {
	newBlog := domainBlog.Blog{
		Title:   blog.Title,
		Content: blog.Content,
		LoginID: blog.LoginID,
	}
	return r.db.Table("BLOGS").Create(&newBlog).Error
}

// ブログを取得
func (r *blogRepository) FindByID(id string) (*domainBlog.Blog, error) {
	blog := domainBlog.Blog{}
	if err := r.db.Table("BLOGS").Where("id = ?", id).First(&blog).Error; err != nil {
		return nil, err
	}
	return &domainBlog.Blog{
		ID:      blog.ID,
		Title:   blog.Title,
		Content: blog.Content,
		LoginID: blog.LoginID,
	}, nil
}

// ユーザーIDに紐づくブログを取得
func (r *blogRepository) FindByUserID(userID string) ([]domainBlog.Blog, error) {
	var blogs []domainBlog.Blog
	if err := r.db.Table("BLOGS").Where("user_id = ?", userID).Find(&blogs).Error; err != nil {
		return nil, err
	}

	var domainBlogs []domainBlog.Blog
	for _, blog := range blogs {
		domainBlogs = append(domainBlogs, domainBlog.Blog{
			ID:      blog.ID,
			Title:   blog.Title,
			Content: blog.Content,
			LoginID: blog.LoginID,
		})
	}
	return domainBlogs, nil
}

// ブログを更新
func (r *blogRepository) Update(blog *domainBlog.Blog) error {
	existingBlog := domainBlog.Blog{}
	if err := r.db.Table("BLOGS").Where("id = ?", blog.ID).First(&existingBlog).Error; err != nil {
		return err
	}

	existingBlog.Title = blog.Title
	existingBlog.Content = blog.Content
	return r.db.Table("BLOGS").Save(&existingBlog).Error
}

// ブログを削除
func (r *blogRepository) Delete(id string) error {
	return r.db.Table("BLOGS").Where("id = ?", id).Delete(&domainBlog.Blog{}).Error
}
