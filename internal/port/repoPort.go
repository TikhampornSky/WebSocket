package port

import (
	"context"
	"server/internal/domain"
)

type UserRepoPort interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	DeleteUserAll(ctx context.Context) error
	UpdateUsername(ctx context.Context, id int64, username string) error
	GetAllUsers(ctx context.Context) ([]*domain.PublicUser, error)
}

type ChatroomRepoPort interface {
	CreateChatroom(ctx context.Context, chatroom *domain.Chatroom) (*domain.Chatroom, error)
	JoinChatroom(ctx context.Context, id int64, clientID int64) (*domain.Chatroom, error)
	LeaveChatroom(ctx context.Context, id int64, clientID int64) error
	GetChatroomByID(ctx context.Context, roomId int64) (*domain.GetRoomByIDRepo, error)
	UpdateChatroomName(ctx context.Context, id int64, name string) error
	GetAllChatrooms(ctx context.Context) ([]*domain.Chatroom, error)
	DeleteChatroomAll(ctx context.Context) error
}
