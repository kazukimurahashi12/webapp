package domain

// encrypting comparing passwordsインターフェース
type Crypto interface {
	Encrypt(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}
