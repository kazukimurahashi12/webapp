package comment

import (
	"time"

	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
)

type Comment struct {
	ID          uint `gorm:"primaryKey"`
	PostID      uint `gorm:"not null"`
	UserID      *uint
	Content     string    `gorm:"type:text;not null"`
	AuthorName  string    `gorm:"size:100"`
	AuthorEmail string    `gorm:"size:100"`
	CreatedAt   time.Time `gorm:"not null"`
	ParentID    *uint
	Status      string          `gorm:"size:20;not null;default:'pending'"`
	User        domainUser.User `gorm:"foreignKey:UserID"`
	Parent      *Comment        `gorm:"foreignKey:ParentID"`
	Replies     []Comment       `gorm:"foreignKey:ParentID"`
}
