package domain

import "time"

type Blog struct {
	Id        uint       `json:"id" gorm:"primarykey"`
	CreatedAt time.Time  `json:"createdAt" gorm:"primarykey"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"primarykey"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index"`
	LoginId   string     `json:"loginID" binding:"required,max=20"`
	Title     string     `json:"title" binding:"required,min=1,max=50"`
	Content   string     `json:"content" binding:"required,min=1,max=8000"`
}

type BlogPost struct {
	Id      string `json:"id"`
	LoginId string `json:"loginID" binding:"required,min=2,max=10"`
	Title   string `json:"title" binding:"required,min=1,max=50"`
	Content string `json:"content" binding:"required,min=1,max=8000"`
}
