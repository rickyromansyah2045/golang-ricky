package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(dbPass, userPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(userPass))

	return err == nil
}
