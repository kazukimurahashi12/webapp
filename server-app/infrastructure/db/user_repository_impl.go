package db

import (
	"github.com/kazukimurahashi12/webapp/crypto"
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/interface/repository"
	"github.com/kazukimurahashi12/webapp/model/entity"
)

type userRepository struct {
	db *DB
}

func NewUserRepository(db *DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	user := entity.User{}
	if err := r.db.Table("USERS").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &domain.User{
		ID:       user.UserId,
		Password: user.Password,
	}, nil
}

func (r *userRepository) FindByUserID(userID string) (*domain.User, error) {
	user := entity.User{}
	if err := r.db.Table("USERS").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &domain.User{
		ID:       user.UserId,
		Password: user.Password,
	}, nil
}

func (r *userRepository) Create(user *domain.User) error {
	encryptPw, err := crypto.PasswordEncrypt(user.Password)
	if err != nil {
		return err
	}

	newUser := entity.User{
		UserId:   user.ID,
		Password: encryptPw,
	}

	return r.db.Table("USERS").Create(&newUser).Error
}

func (r *userRepository) Update(user *domain.User) error {
	existingUser := entity.User{}
	if err := r.db.Table("USERS").Where("user_id = ?", user.ID).First(&existingUser).Error; err != nil {
		return err
	}

	existingUser.Password = user.Password
	return r.db.Table("USERS").Save(&existingUser).Error
}

func (r *userRepository) UpdateID(oldID, newID string) (*domain.User, error) {
	user := entity.User{}
	if err := r.db.Table("USERS").Where("user_id = ?", oldID).First(&user).Error; err != nil {
		return nil, err
	}

	user.UserId = newID
	if err := r.db.Table("USERS").Save(&user).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       user.UserId,
		Password: user.Password,
	}, nil
}

func (r *userRepository) UpdatePassword(userID, newPassword string) (*domain.User, error) {
	user := entity.User{}
	if err := r.db.Table("USERS").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	encryptPw, err := crypto.PasswordEncrypt(newPassword)
	if err != nil {
		return nil, err
	}

	user.Password = encryptPw
	if err := r.db.Table("USERS").Save(&user).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       user.UserId,
		Password: user.Password,
	}, nil
}
