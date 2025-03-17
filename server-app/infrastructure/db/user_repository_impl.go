package db

import (
	"github.com/kazukimurahashi12/webapp/domain"
	"github.com/kazukimurahashi12/webapp/infrastructure/crypto"
	"github.com/kazukimurahashi12/webapp/interface/repository"
)

type userRepository struct {
	db *DB
}

func NewUserRepository(db *DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	user := domain.User{}
	if err := r.db.Table("USERS").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &domain.User{
		Id:       user.Id,
		Password: user.Password,
	}, nil
}

func (r *userRepository) FindByUserID(userID string) (*domain.User, error) {
	user := domain.User{}
	if err := r.db.Table("USERS").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &domain.User{
		UserId:   user.UserId,
		Password: user.Password,
	}, nil
}

func (r *userRepository) Create(user *domain.User) error {
	crypto := crypto.NewBcryptCrypto()
	encryptPw, err := crypto.Encrypt(user.Password)
	if err != nil {
		return err
	}

	newUser := domain.User{
		UserId:   user.UserId,
		Password: encryptPw,
	}

	return r.db.Table("USERS").Create(&newUser).Error
}

func (r *userRepository) Update(user *domain.User) error {
	existingUser := domain.User{}
	if err := r.db.Table("USERS").Where("user_id = ?", user.Id).First(&existingUser).Error; err != nil {
		return err
	}

	existingUser.Password = user.Password
	return r.db.Table("USERS").Save(&existingUser).Error
}

func (r *userRepository) UpdateID(oldID, newID string) (*domain.User, error) {
	user := domain.User{}
	if err := r.db.Table("USERS").Where("user_id = ?", oldID).First(&user).Error; err != nil {
		return nil, err
	}

	user.UserId = newID
	if err := r.db.Table("USERS").Save(&user).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		UserId:   user.UserId,
		Password: user.Password,
	}, nil
}

func (r *userRepository) UpdatePassword(userID, newPassword string) (*domain.User, error) {
	user := domain.User{}
	if err := r.db.Table("USERS").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	crypto := crypto.NewBcryptCrypto()
	encryptPw, err := crypto.Encrypt(newPassword)
	if err != nil {
		return nil, err
	}

	user.Password = encryptPw
	if err := r.db.Table("USERS").Save(&user).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		UserId:   user.UserId,
		Password: user.Password,
	}, nil
}
