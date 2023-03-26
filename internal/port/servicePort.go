package port

import (
	"context"
	"server/internal/domain"
)

type UserServicePort interface {
	CreateUser(ctx context.Context, req *domain.CreateUserReq) (*domain.CreateUserRes, error)
	Login(c context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error)
	UpdateUsername(ctx context.Context, req *domain.UpdateUsernameReq) error
	GetAllUsers(ctx context.Context) ([]*domain.PublicUser, error)
}
