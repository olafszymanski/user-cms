package users

import (
	"fmt"
	"strings"

	"github.com/olafszymanski/user-cms/postgres"
	"github.com/olafszymanski/user-cms/utils"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

func (u *User) generateUpdateQuery(username, email, password *string, admin *bool) (string, error) {
	var query strings.Builder
	query.WriteString("UPDATE users SET ")

	params := []string{}
	if username != nil {
		params = append(params, "username=:username")
		u.Username = *username
	}
	if email != nil {
		params = append(params, "email=:email")
		u.Email = *email
	}
	if password != nil {
		params = append(params, "password=:password")
		hash, err := utils.HashPassword(*password)
		if err != nil {
			return "", fmt.Errorf("could not hash password, error: %w", err)
		}
		u.Password = hash
	}
	if admin != nil {
		params = append(params, "admin=:admin")
		u.Admin = *admin
	}

	for i, p := range params {
		query.WriteString(p)
		if i != len(params)-1 {
			query.WriteByte(',')
		}
	}
	query.WriteString(" WHERE id=:id")
	return query.String(), nil

	// query := "UPDATE users SET "
	// params := []string{}
	// if username != nil {
	// 	params = append(params, "username=:username")
	// 	u.Username = *username
	// }
	// if email != nil {
	// 	params = append(params, "email=:email")
	// 	u.Email = *email
	// }
	// if password != nil {
	// 	params = append(params, "password=:password")
	// 	hash, err := utils.HashPassword(*password)
	// 	if err != nil {
	// 		return "", fmt.Errorf("could not hash password, error: %w", err)
	// 	}
	// 	u.Password = hash
	// }
	// if admin != nil {
	// 	params = append(params, "admin=:admin")
	// 	u.Admin = *admin
	// }

	// for i, p := range params {
	// 	query += p
	// 	if i != len(params)-1 {
	// 		query += ", "
	// 	}
	// }

	// query += " WHERE id=:id"
	// return query, nil
}

func (u *User) Create() error {
	password, err := utils.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("could not hash password, error: %w", err)
	}

	id := 0
	if err := postgres.Database.QueryRow("INSERT INTO users (username, email, password, admin) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Username,
		u.Email,
		string(password),
		utils.Btou(u.Admin)).Scan(&id); err != nil {
		return fmt.Errorf("could not insert user into database, error: %w", err)
	}
	u.ID = int(id)
	return nil
}

func (u *User) Update(username, password, email *string, admin *bool) error {
	if err := postgres.Database.QueryRowx("SELECT * FROM users WHERE id = $1", u.ID).StructScan(u); err != nil {
		return fmt.Errorf("user with id %v does not exist, error: %w", u.ID, err)
	}

	query, err := u.generateUpdateQuery(username, password, email, admin)
	if err != nil {
		return fmt.Errorf("could not generate update query, error: %w", err)
	}

	if _, err := postgres.Database.NamedQuery(query, map[string]interface{}{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
		"admin":    utils.Btou(u.Admin),
	}); err != nil {
		return fmt.Errorf("could not update user, error: %w", err)
	}
	return nil
}

func (u *User) Delete() error {
	if err := postgres.Database.QueryRowx("SELECT * FROM users WHERE id = $1", u.ID).StructScan(&User{}); err != nil {
		return fmt.Errorf("user with id %v does not exist, error: %w", u.ID, err)
	}
	if _, err := postgres.Database.Exec("DELETE FROM users WHERE id = $1", u.ID); err != nil {
		return fmt.Errorf("could not delete user with id %v, error: %w", u.ID, err)
	}
	return nil
}

func Get(id int) (*User, error) {
	user := &User{}
	if err := postgres.Database.QueryRowx("SELECT * FROM users WHERE id = $1", id).StructScan(user); err != nil {
		return nil, fmt.Errorf("user with id %v does not exist, error: %w", id, err)
	}
	return user, nil
}

func All() ([]*User, error) {
	var users []*User
	if err := postgres.Database.Select(&users, "SELECT * FROM users"); err != nil {
		return nil, nil
	}
	return users, nil
}
