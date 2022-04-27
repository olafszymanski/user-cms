package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Btou(v bool) uint8 {
	if v {
		return 1
	} else {
		return 0
	}
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
