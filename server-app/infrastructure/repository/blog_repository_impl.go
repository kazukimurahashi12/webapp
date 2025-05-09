package repository

import (
	"errors"
	"fmt"

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
	return r.db.Table("BLOGS").Create(blog).Error
}

// ブログを取得
func (r *blogRepository) FindBlogByID(id uint) (*domainBlog.Blog, error) {
	blog := domainBlog.Blog{}
	if err := r.db.Table("BLOGS").Where("id = ?", id).First(&blog).Error; err != nil {
		return nil, fmt.Errorf("failed to find blog (id=%d): %w", id, err)
	}
	return &blog, nil
}

// 著者IDに紐づくブログを取得
func (r *blogRepository) FindBlogsByAuthorID(authorID uint) ([]domainBlog.Blog, error) {
	var blogs []domainBlog.Blog
	if err := r.db.Table("BLOGS").Where("author_id = ?", authorID).Find(&blogs).Error; err != nil {
		return nil, fmt.Errorf("failed to find blogs by author_id (author_id=%d): %w", authorID, err)
	}
	return blogs, nil
}

// 著者IDに対応するブログを取得
func (r *blogRepository) FindBlogByAuthorID(authorID uint) (*domainBlog.Blog, error) {
	blog := domainBlog.Blog{}
	if err := r.db.Table("BLOGS").Where("author_id = ?", authorID).First(&blog).Error; err != nil {
		return nil, fmt.Errorf("failed to find blog by author_id (author_id=%d): %w", authorID, err)
	}
	return &blog, nil
}

// ブログを更新
func (r *blogRepository) Update(blog *domainBlog.Blog) (err error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	existingBlog := domainBlog.Blog{}
	if err = tx.Table("BLOGS").Where("id = ?", blog.ID).First(&existingBlog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domainBlog.ErrBlogNotFound
		}
		return fmt.Errorf("failed to find existing blog (id=%d): %w", blog.ID, err)
	}

	updateData := map[string]interface{}{
		"title":   blog.Title,
		"content": blog.Content,
	}

	if err = tx.Table("BLOGS").Where("id = ?", blog.ID).Updates(updateData).Error; err != nil {
		return fmt.Errorf("failed to update blog (id=%d): %w", blog.ID, err)
	}

	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ブログを削除
func (r *blogRepository) Delete(id uint) error {
	if err := r.db.Table("BLOGS").Where("id = ?", id).Delete(&domainBlog.Blog{}).Error; err != nil {
		return fmt.Errorf("failed to delete blog (id=%d): %w", id, err)
	}
	return nil
}
