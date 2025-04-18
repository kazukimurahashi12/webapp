package dto

import (
	"time"
)

type Blog struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
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

type BlogPostResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Created string `json:"created_at"`
}
