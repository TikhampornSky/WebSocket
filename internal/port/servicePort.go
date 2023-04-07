package port

import (
	"context"
	"server/internal/domain"
)

type UserServicePort interface {
	CreateUser(ctx context.Context, req *domain.CreateUserReq) (*domain.CreateUserRes, error)
	Login(c context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error)
	UpdateUser(ctx context.Context, req *domain.UpdateUsernameReq) error
	UpdatePassword(ctx context.Context, req *domain.UpdatePasswordReq) error
	GetAllUsers(ctx context.Context) ([]*domain.PublicUser, error)
}

type ChatroomServicePort interface {
	CreateChatroom(ctx context.Context, req *domain.CreateChatroomReq) (*domain.CreateChatroomRes, error)
	JoinChatroom(ctx context.Context, req *domain.JoinLeaveChatroomReq) (*domain.JoinLeaveChatroomRes, error)
	LeaveChatroom(ctx context.Context, req *domain.JoinLeaveChatroomReq) error
	GetChatroomByID(ctx context.Context, req *domain.GetChatroomByIDReq) (*domain.GetChatroomByIDRes, error)
	UpdateChatroomName(ctx context.Context, req *domain.UpdateChatroomNameReq) error
	GetAllChatrooms(ctx context.Context) ([]*domain.Chatroom, error)
}
