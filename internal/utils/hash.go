package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultCost = 12
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func VerifyPassword(password, hash string) bool {
	return CheckPassword(password, hash)
}
