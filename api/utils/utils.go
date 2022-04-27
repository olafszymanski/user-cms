package utils

import (
	"fmt"

	"github.com/olafszymanski/user-cms/graph/model"
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

func GenerateQueryAndUser(input *model.UpdateUser, user *model.User) (string, error) {
	query := "UPDATE users SET "

	params := []string{}
	if input.Username != nil {
		params = append(params, "username=:username")
		user.Username = input.Username
	}
	if input.Email != nil {
		params = append(params, "email=:email")
		user.Email = input.Email
	}
	if input.Password != nil {
		params = append(params, "password=:password")
		password, err := HashPassword(*input.Password)
		if err != nil {
			return "", fmt.Errorf("could not hash password, error: %w", err)
		}
		input.Password = &password
	}
	if input.Admin != nil {
		params = append(params, "admin=:admin")
		user.Admin = input.Admin
	}

	for i, p := range params {
		query += p
		if i != len(params)-1 {
			query += ", "
		}
	}

	query += " WHERE id=:id"
	return query, nil
}
