package utils

import (
	"github.com/olafszymanski/user-cms/graph/model"
)

func Btou(v bool) uint8 {
	if v {
		return 1
	} else {
		return 0
	}
}

func BuildUpdateQuery(user *model.UpdateUser) string {
	query := "UPDATE users SET "
	params := []string{}
	if user.Username != nil {
		params = append(params, "username=:username")
	}
	if user.Email != nil {
		params = append(params, "email=:email")
	}
	if user.Password != nil {
		params = append(params, "password=:password")
	}
	if user.Admin != nil {
		params = append(params, "admin=:admin")
	}

	for i, p := range params {
		query += p
		if i != len(params)-1 {
			query += ", "
		}
	}

	query += " WHERE id=:id"
	return query
}
