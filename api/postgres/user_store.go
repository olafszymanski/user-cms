package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olafszymanski/user-cms/users"
	"github.com/olafszymanski/user-cms/utils"
)

func NewUserStore(db *sqlx.DB) *UserStore {
	store := &UserStore{
		DB:         db,
		statements: make(map[string]*sqlx.Stmt, 0),
	}

	// Prepared statements
	queries := []string{
		"SELECT * FROM users WHERE id = $1",
		"SELECT * FROM users",
		"INSERT INTO users (username, email, password, admin) VALUES ($1, $2, $3, $4) RETURNING id",
		"DELETE FROM users WHERE id = $1",
		"SELECT * FROM users WHERE username = $1",
	}
	for _, query := range queries {
		stmt, err := store.Preparex(query)
		if err != nil {
			panic(fmt.Errorf("could not prepare query: %v, error: %w", query, err))
		}
		store.statements[query] = stmt
	}

	return store
}

type UserStore struct {
	*sqlx.DB
	statements map[string]*sqlx.Stmt
}

func (u *UserStore) Get(id int) (*users.User, error) {
	user := &users.User{}
	if err := u.statements["SELECT * FROM users WHERE id = $1"].QueryRowx(id).StructScan(user); err != nil {
		return nil, fmt.Errorf("user with id %v does not exist, error: %w", id, err)
	}
	return user, nil
}

func (u *UserStore) All() ([]*users.User, error) {
	var users []*users.User
	if err := u.statements["SELECT * FROM users"].Select(&users); err != nil {
		return nil, nil
	}
	return users, nil
}

func (u *UserStore) Create(new *users.User) (*users.User, error) {
	password, err := utils.HashPassword(*new.Password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password, error: %w", err)
	}

	id := 0
	if err := u.statements["INSERT INTO users (username, email, password, admin) VALUES ($1, $2, $3, $4) RETURNING id"].QueryRow(
		new.Username,
		new.Email,
		password,
		utils.Btou(*new.Admin)).Scan(&id); err != nil {
		return nil, fmt.Errorf("could not insert user into database, error: %w", err)
	}
	new.ID = &id
	new.Password = &password
	return new, nil
}

func (u *UserStore) Update(update *users.User) (*users.User, error) {
	user := &users.User{}
	if err := u.statements["SELECT * FROM users WHERE id = $1"].QueryRowx(update.ID).StructScan(user); err != nil {
		return nil, fmt.Errorf("user with id %v does not exist, error: %w", update.ID, err)
	}

	query, err := utils.GenerateUpdateQuery(user, update)
	if err != nil {
		return nil, fmt.Errorf("could not generate update query, error: %w", err)
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

func (u *UserStore) Delete(id int) error {
	if err := u.statements["SELECT * FROM users WHERE id = $1"].QueryRowx(id).StructScan(&users.User{}); err != nil {
		return fmt.Errorf("user with id %v does not exist, error: %w", id, err)
	}
	if _, err := u.Exec("DELETE FROM users WHERE id = $1", id); err != nil {
		return fmt.Errorf("could not delete user with id %v, error: %w", id, err)
	}
	return nil
}

func (u *UserStore) GetByUsername(username string) (*users.User, error) {
	user := &users.User{}
	if err := u.statements["SELECT * FROM users WHERE username = $1"].QueryRowx(username).StructScan(user); err != nil {
		return nil, fmt.Errorf("user with username '%v' does not exist, error: %w", username, err)
	}
	return user, nil
}
