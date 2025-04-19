package crypto

import (
	"github.com/kazukimurahashi12/webapp/domain/crypto"
	"golang.org/x/crypto/bcrypt"
)

type BcryptCrypto struct{}

func NewBcryptCrypto() crypto.Crypto {
	return &BcryptCrypto{}
}

// bcryptを使用してパスワードをハッシュ化
func (b *BcryptCrypto) Encrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// ハッシュ化されたパスワードと入力されたパスワードが一致するかどうかを確認
func (b *BcryptCrypto) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
