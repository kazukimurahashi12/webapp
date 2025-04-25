package user

type UserRepository interface {
	FindUserByID(id uint) (*User, error)
	FindUserByUserID(userID string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	UpdateID(oldID, newID string) (*User, error)
	UpdatePassword(userID, newPassword string) (*User, error)
}
