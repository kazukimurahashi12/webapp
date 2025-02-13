package entity

import (
	"time"
)

type Blog struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	CreatedAt time.Time  `json:"createdAt" gorm:"primarykey"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"primarykey"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index"`
	LoginID   string     `json:"loginID" binding:"required,max=20"`
	Title     string     `json:"title" binding:"required,min=1,max=50"`
	Content   string     `json:"content" binding:"required,min=1,max=8000"`
}

type BlogPost struct {
	ID      string `json:"id"`
	LoginID string `json:"loginID" binding:"required,min=2,max=10"`
	Title   string `json:"title" binding:"required,min=1,max=50"`
	Content string `json:"content" binding:"required,min=1,max=8000"`
}
