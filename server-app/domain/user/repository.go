package user

type UserRepository interface {
	FindUserByID(id uint) (*User, error)
	FindUserByUserID(userID uint) (*User, error)
	Create(user *User) error
	Update(user *User) error
	UpdateID(oldID, newID uint) (*User, error)
	UpdatePassword(userID uint, newPassword string) (*User, error)
}
