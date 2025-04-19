package repository

import (
	"errors"

	domainUser "github.com/kazukimurahashi12/webapp/domain/user"
	"github.com/kazukimurahashi12/webapp/infrastructure/crypto"
	"github.com/kazukimurahashi12/webapp/infrastructure/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(manager *db.DBManager) domainUser.UserRepository {
	return &userRepository{
		db:     manager.DB,
		logger: manager.Logger}
}

// シーケンシャルIDによりユーザーを取得
func (r *userRepository) FindByID(id string) (*domainUser.User, error) {
	user := domainUser.User{}
	if err := r.db.Table("USERS").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &domainUser.User{
		ID:       user.ID,
		Password: user.Password,
	}, nil
}

// ユーザーIDに紐づくユーザーを取得
func (r *userRepository) FindByUserID(userID string) (*domainUser.User, error) {
	user := domainUser.User{}
	if err := r.db.Table("USERS").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &domainUser.User{
		UserID:   user.UserID,
		Password: user.Password,
	}, nil
}

// ユーザーを作成
func (r *userRepository) Create(user *domainUser.User) error {
	crypto := crypto.NewBcryptCrypto()
	encryptPw, err := crypto.Encrypt(user.Password)
	if err != nil {
		return err
	}

	newUser := domainUser.User{
		UserID:   user.UserID,
		Password: encryptPw,
	}

	return r.db.Table("USERS").Create(&newUser).Error
}

// ユーザーIDを変更
func (r *userRepository) UpdateID(oldID, newID string) (*domainUser.User, error) {
	// トランザクション開始
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := domainUser.User{}
	if err := tx.Table("USERS").Where("user_id = ?", oldID).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainUser.ErrUserNotFound
		}
		return nil, err
	}

	user.UserID = newID
	if err := tx.Table("USERS").Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	result := &domainUser.User{
		UserID:   user.UserID,
		Password: user.Password,
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return result, nil
}

// ユーザーPWを変更
func (r *userRepository) UpdatePassword(userID, newPassword string) (*domainUser.User, error) {
	// トランザクション開始
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := domainUser.User{}
	if err := tx.Table("USERS").Where("user_id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainUser.ErrUserNotFound
		}
		return nil, err
	}

	crypto := crypto.NewBcryptCrypto()
	encryptPw, err := crypto.Encrypt(newPassword)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	user.Password = encryptPw
	if err := tx.Table("USERS").Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	result := &domainUser.User{
		UserID:   user.UserID,
		Password: user.Password,
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return result, nil
}

// ユーザー情報を更新
func (r *userRepository) Update(user *domainUser.User) error {
	// トランザクション開始
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existingUser := domainUser.User{}
	if err := tx.Table("USERS").Where("user_id = ?", user.ID).First(&existingUser).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domainUser.ErrUserNotFound
		}
		return err
	}

	existingUser.Password = user.Password
	if err := tx.Table("USERS").Save(&existingUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
