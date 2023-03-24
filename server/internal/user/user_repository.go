package user

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	queryEmail := "SELECT id FROM users WHERE email = $1"
	var idEmail int64
	err := r.db.QueryRowContext(ctx, queryEmail, user.Email).Scan(&idEmail)
	if err == nil {
		return &User{}, fmt.Errorf("user with email %s already exists", user.Email)
	}

	queryUsername := "SELECT id FROM users WHERE username = $1"
	var idUsername int64
	err = r.db.QueryRowContext(ctx, queryUsername, user.Username).Scan(&idUsername)
	if err == nil {
		return &User{}, fmt.Errorf("user with username %s already exists", user.Username)
	}

	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int64
	err = r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return &User{}, err
	}

	user.ID = id
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}
	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &User{}, nil
	}

	return &u, nil
}
