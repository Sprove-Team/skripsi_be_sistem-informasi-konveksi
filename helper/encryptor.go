package helper

import "golang.org/x/crypto/bcrypt"

type Encryptor interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type encryptor struct{}

func NewEncryptor() Encryptor {
	return &encryptor{}
}

func (e *encryptor) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (e *encryptor) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
