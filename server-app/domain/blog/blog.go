package blog

import (
	"time"
)

type Blog struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	LoginID   string     `json:"loginID" binding:"required,max=20"`
	Title     string     `json:"title" binding:"required,min=1,max=50"`
	Content   string     `json:"content" binding:"required,min=1,max=8000"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index"`
}

type BlogPost struct {
	ID      string `json:"id" gorm:"primaryKey"`
	LoginID string `json:"loginID" binding:"required,min=2,max=10"`
	Title   string `json:"title" binding:"required,min=1,max=50"`
	Content string `json:"content" binding:"required,min=1,max=8000"`
}
