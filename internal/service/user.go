package service

import (
	"context"
	"server/internal/domain"
	"server/internal/port"
	"server/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey = "secret"
)

type userService struct {
	port.UserRepoPort
	timeout time.Duration
}

func NewUserService(repo port.UserRepoPort) port.UserServicePort {
	return &userService{
		repo,
		time.Duration(2) * time.Second,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *domain.CreateUserReq) (*domain.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.UserRepoPort.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &domain.CreateUserRes{
		ID:       r.ID,
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *userService) Login(c context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	u, err := s.UserRepoPort.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &domain.LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &domain.LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &domain.LoginUserRes{}, err
	}

	return &domain.LoginUserRes{
		AccessToken: ss, 
		Username: u.Username, 
		ID: strconv.Itoa(int(u.ID)),
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *domain.UpdateUsernameReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	id := req.ID
	err := s.UserRepoPort.UpdateUser(ctx, id, req.Username, req.Email)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) UpdatePassword(ctx context.Context, req *domain.UpdatePasswordReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	id := req.ID
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}
	err = s.UserRepoPort.UpdatePassword(ctx, id, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*domain.PublicUser, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	users, err := s.UserRepoPort.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
