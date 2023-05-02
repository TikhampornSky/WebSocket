package repo

import (
	"context"
	"database/sql"
	"server/internal/domain"
	"server/internal/port"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type userRepository struct {
	db DBTX
}

func NewUserRepository(db DBTX) port.UserRepoPort {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	queryEmail := "SELECT id FROM users WHERE email = $1"
	var idEmail int64
	err := r.db.QueryRowContext(ctx, queryEmail, user.Email).Scan(&idEmail)
	if err == nil {
		return &domain.User{}, domain.ErrDuplicateEmail.With("user with email %s already exists", user.Email)
	}

	queryUsername := "SELECT id FROM users WHERE username = $1"
	var idUsername int64
	err = r.db.QueryRowContext(ctx, queryUsername, user.Username).Scan(&idUsername)
	if err == nil {
		return &domain.User{}, domain.ErrDuplicateUsername.With("user with username %s already exists", user.Username)
	}

	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int64
	err = r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return &domain.User{}, domain.ErrInternal.From(err.Error(), err)
	}

	user.ID = id
	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := domain.User{}
	query := "SELECT id, email, username, password FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	if err != nil {
		return &domain.User{}, nil
	}

	return &u, nil
}

func (r *userRepository) DeleteUserAll(ctx context.Context) error { // Testing Propose
	query := "DELETE FROM users WHERE id > 0"
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}

	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id int64, username, email string) error {
	queryFindOldUsername := "SELECT username FROM users WHERE id = $1"
	var oldUsername string
	err := r.db.QueryRowContext(ctx, queryFindOldUsername, id).Scan(&oldUsername)
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}
	query := "UPDATE users SET username = $1, email = $2 WHERE id = $3"
	_, err = r.db.ExecContext(ctx, query, username, email, id)
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}

	// update name of room type "private"
	query = "UPDATE chatrooms SET name = REPLACE(name, $1, $2) WHERE category = 'private' AND $3 = ANY (clients)"
	_, err = r.db.ExecContext(ctx, query, oldUsername, username, id)
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}

	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	query := "UPDATE users SET password = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, password, id)
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}
	return nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]*domain.PublicUser, error) {
	query := "SELECT id, email, username FROM users"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, domain.ErrInternal.From(err.Error(), err)
	}

	var users []*domain.PublicUser
	for rows.Next() {
		u := domain.PublicUser{}
		err := rows.Scan(&u.ID, &u.Email, &u.Username)
		if err != nil {
			return nil, domain.ErrInternal.From(err.Error(), err)
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *userRepository) DeleteAllUsers(ctx context.Context) error {
	query := "DELETE FROM users"
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}

	return nil
}
