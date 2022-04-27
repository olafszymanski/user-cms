package utils

import (
	"fmt"
	"strings"

	"github.com/olafszymanski/user-cms/users"
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

func GenerateUpdateQuery(before, after *users.User) (string, error) {
	var query strings.Builder
	query.WriteString("UPDATE users SET ")

	params := []string{}
	if after.Username != nil {
		params = append(params, "username=:username")
		before.Username = after.Username
	}
	if after.Email != nil {
		params = append(params, "email=:email")
		before.Email = after.Email
	}
	if after.Password != nil {
		params = append(params, "password=:password")
		hash, err := HashPassword(*after.Password)
		if err != nil {
			return "", fmt.Errorf("could not hash password, error: %w", err)
		}
		before.Password = &hash
	}
	if after.Admin != nil {
		params = append(params, "admin=:admin")
		before.Admin = after.Admin
	}

	for i, p := range params {
		query.WriteString(p)
		if i != len(params)-1 {
			query.WriteByte(',')
		}
	}
	query.WriteString(" WHERE id=:id")
	return query.String(), nil
}
