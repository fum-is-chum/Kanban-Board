package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func VerifyPassword(hashedPass string, plainText string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainText))
	if err != nil {
		return errors.New("Password Invalid")
	}
	return nil
}
