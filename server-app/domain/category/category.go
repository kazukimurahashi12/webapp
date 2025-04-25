package category

import (
	domainBlog "github.com/kazukimurahashi12/webapp/domain/blog"
)

type Category struct {
	ID          uint              `gorm:"primaryKey"`
	Name        string            `gorm:"size:50;not null;unique"`
	Description string            `gorm:"size:255"`
	ParentID    *uint             `gorm:"default:null"`
	Parent      *Category         `gorm:"foreignKey:ParentID"`
	Children    []Category        `gorm:"foreignKey:ParentID"`
	Blog        []domainBlog.Blog `gorm:"many2many:post_categories;"`
}
