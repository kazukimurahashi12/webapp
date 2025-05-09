package blog

import (
	"time"

	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
)

type Blog struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	AuthorID  uint            `json:"authorId" gorm:"column:user_id"`                  // 外部キー
	Author    domainUser.User `json:"author" gorm:"foreignKey:AuthorID;references:ID"` // Userへ参照
	Title     string          `json:"title" binding:"required,min=1,max=50"`
	Content   string          `json:"content" binding:"required,min=1,max=8000"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt *time.Time      `json:"deletedAt" gorm:"index"`
}
