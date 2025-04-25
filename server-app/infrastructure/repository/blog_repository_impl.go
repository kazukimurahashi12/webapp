package repository

import (
	"errors"

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
func (r *blogRepository) Create(blog *domainBlog.Blog) error {
	return r.db.Table("BLOGS").Create(&blog).Error
}

// ブログを取得
func (r *blogRepository) FindBlogByID(id uint) (*domainBlog.Blog, error) {
	blog := domainBlog.Blog{}
	if err := r.db.Table("BLOGS").Where("id = ?", id).First(&blog).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

// ユーザーIDに紐づくブログを取得
func (r *blogRepository) FindBlogsByUserID(userID string) ([]domainBlog.Blog, error) {
	var blogs []domainBlog.Blog
	if err := r.db.Table("BLOGS").Where("user_id = ?", userID).Find(&blogs).Error; err != nil {
		return nil, err
	}

	return blogs, nil
}

// ブログを更新
func (r *blogRepository) Update(blog *domainBlog.Blog) error {
	// トランザクション開始
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existingBlog := domainBlog.Blog{}
	if err := tx.Table("BLOGS").Where("id = ?", blog.ID).First(&existingBlog).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domainBlog.ErrBlogNotFound
		}
		return err
	}

	// 更新対象フィールドを明示的に指定
	updateData := map[string]interface{}{
		"title":   blog.Title,
		"content": blog.Content,
	}

	if err := tx.Table("BLOGS").Where("id = ?", blog.ID).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return err
	}

	// コミット
	return tx.Commit().Error
}

// ブログを削除
func (r *blogRepository) Delete(id string) error {
	return r.db.Table("BLOGS").Where("id = ?", id).Delete(&domainBlog.Blog{}).Error
}
