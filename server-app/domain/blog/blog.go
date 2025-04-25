package blog

import (
	"time"

	domaiUser "github.com/kazukimurahashi12/webapp/domain/user"
)

type Blog struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null"`
	User      domaiUser.User `gorm:"foreignKey:UserID"`
	Title     string         `binding:"required,min=1,max=50"`
	Content   string         `binding:"required,min=1,max=8000"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt *time.Time     `json:"deletedAt" gorm:"index"`
}
