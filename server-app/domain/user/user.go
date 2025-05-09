package user

import "time"

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Username  string     `json:"username" gorm:"column:user_id;uniqueIndex" binding:"required,min=2,max=10"`
	Password  string     `json:"password" binding:"required,min=4,max=20"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index"`
}
