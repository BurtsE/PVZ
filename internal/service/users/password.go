package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func checkPassword(password string) error {
	if len(password) < 8 || len(password) > 15 {
		return errors.New("Password must be between 8 and 15 characters")
	}
	return nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func compareHashAndPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
