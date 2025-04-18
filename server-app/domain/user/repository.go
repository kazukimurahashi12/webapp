package user

type UserRepository interface {
	FindByID(id string) (*User, error)
	FindByUserID(userID string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	UpdateID(oldID, newID string) (*User, error)
	UpdatePassword(userID, newPassword string) (*User, error)
}
