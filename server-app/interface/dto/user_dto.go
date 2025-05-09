package dto

import "time"

type User struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	UserID    string     `json:"userId" binding:"required,min=2,max=10"`
	Password  string     `json:"password" binding:"required,min=4,max=20"`
}

type UserIdChange struct {
	CurrentId string `json:"currentId" binding:"required,min=2,max=10"`
	NewId     string `json:"newId" binding:"required,min=2,max=10"`
}

type UserPwChange struct {
	UserID         string `json:"userId" binding:"required,min=2,max=10"`
	NowPassword    string `json:"nowPassword" binding:"required,min=4,max=20"`
	ChangePassword string `json:"changePassword" binding:"required,min=4,max=20"`
}

type FormUser struct {
	UserID   string `json:"userId" binding:"required,min=2,max=10"`
	Password string `json:"password" binding:"required,min=4,max=20"`
}

type UserCreatedResponse struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"userId" binding:"required,min=2,max=10"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserIDResponse struct {
	UserID string `json:"userId" binding:"required,min=2,max=10"`
}
