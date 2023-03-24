package user

import (
	"context"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	DeleteUserAll(ctx context.Context) error
	UpdateUsername(ctx context.Context, id int64, username string) error
	GetAllUsers(ctx context.Context) ([]*PublicUser, error)
}

type Service interface {
	CreateUser(ctx context.Context, req *CreateUserReq) (*CreateUserRes, error)
	Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
	UpdateUsername(ctx context.Context, req *UpdateUsernameReq) error
	GetAllUsers(ctx context.Context) ([]*PublicUser, error)
}

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	accessToken string
	ID          string `json:"id"`
	Username    string `json:"username"`
}

type UpdateUsernameReq struct {
	ID       string  `json:"id"`
	Username string `json:"username"`
}

type PublicUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
