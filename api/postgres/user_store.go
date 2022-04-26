package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olafszymanski/user-cms/graph/model"
	"github.com/olafszymanski/user-cms/utils"
	"golang.org/x/crypto/bcrypt"
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

func (u *UserStore) CreateUser(user *model.NewUser) (*model.User, error) {
	if err := u.QueryRowx("SELECT * FROM users WHERE username = $1 OR email = $2", user.Username, user.Email).StructScan(&model.User{}); err != nil {
		// User with specified credentials does not exist
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		fmt.Println(len(password))
		if err != nil {
			return nil, fmt.Errorf("could not hash password")
		}
		if _, err := u.Exec("INSERT INTO users (username, email, password, admin) VALUES ($1, $2, $3, $4)",
			user.Username,
			user.Email,
			string(password),
			utils.Btou(user.Admin)); err != nil {
			return nil, fmt.Errorf("could not insert user into database, error: %w", err)
		}
		return &model.User{
			Username: &user.Username,
			Email:    &user.Email,
			Password: &user.Password,
			Admin:    &user.Admin,
		}, nil
	}
	return nil, fmt.Errorf("user with specified credentials already exists")
}

func (u *UserStore) UpdateUser(user *model.UpdateUser) (*model.User, error) {
	if err := u.QueryRowx("SELECT * FROM users WHERE id = $1", user.ID).StructScan(&model.User{}); err != nil {
		return nil, fmt.Errorf("user with id %v does not exist, error: %w", user.ID, err)
	}
	password, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password")
	}
	stringPassword := string(password)
	user.Password = &stringPassword
	if _, err := u.NamedQuery(utils.BuildUpdateQuery(user), user); err != nil {
		return nil, fmt.Errorf("could not update user, error: %w", err)
	}
	return &model.User{
		ID:       &user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Admin:    user.Admin,
	}, nil
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
