package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olafszymanski/user-cms/graph/model"
	"github.com/olafszymanski/user-cms/utils"
)

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

type UserStore struct {
	*sqlx.DB
}

func (u *UserStore) User(id int) (*model.User, error) {
	user := &model.User{}
	if err := u.QueryRowx("SELECT * FROM users WHERE id = $1", id).StructScan(user); err != nil {
		return nil, fmt.Errorf("user with id %v does not exist, error: %w", id, err)
	}
	return user, nil
}

func (u *UserStore) Users() ([]*model.User, error) {
	var users []*model.User
	if err := u.Select(&users, "SELECT * FROM users"); err != nil {
		// No users in the database
		return nil, nil
	}
	return users, nil
}

func (u *UserStore) CreateUser(input *model.NewUser) (*model.User, error) {
	if err := u.QueryRowx("SELECT * FROM users WHERE username = $1 OR email = $2", input.Username, input.Email).StructScan(&model.User{}); err != nil {
		// User with specified credentials does not exist
		password, err := utils.HashPassword(input.Password)
		if err != nil {
			return nil, fmt.Errorf("could not hash password, error: %w", err)
		}

		if _, err := u.Exec("INSERT INTO users (username, email, password, admin) VALUES ($1, $2, $3, $4)",
			input.Username,
			input.Email,
			string(password),
			utils.Btou(input.Admin)); err != nil {
			return nil, fmt.Errorf("could not insert user into database, error: %w", err)
		}

		return &model.User{
			Username: &input.Username,
			Email:    &input.Email,
			Password: &password,
			Admin:    &input.Admin,
		}, nil
	}
	return nil, fmt.Errorf("user with specified credentials already exists")
}

func (u *UserStore) UpdateUser(input *model.UpdateUser) (*model.User, error) {
	user := &model.User{}
	if err := u.QueryRowx("SELECT * FROM users WHERE id = $1", input.ID).StructScan(user); err != nil {
		return nil, fmt.Errorf("user with id %v does not exist, error: %w", input.ID, err)
	}

	query, err := utils.GenerateQueryAndUser(input, user)
	if err != nil {
		return nil, fmt.Errorf("could not generate query and user, error: %w", err)
	}
	if _, err := u.NamedQuery(query, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
		"admin":    utils.Btou(*user.Admin),
	}); err != nil {
		return nil, fmt.Errorf("could not update user, error: %w", err)
	}
	return user, nil
}

func (u *UserStore) DeleteUser(id int) error {
	if err := u.QueryRowx("SELECT * FROM users WHERE id = $1", id).StructScan(&model.User{}); err != nil {
		return fmt.Errorf("user with id %v does not exist, error: %w", id, err)
	}

	if _, err := u.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		return fmt.Errorf("could not delete user with id %v, error: %w", id, err)
	}
	return nil
}
