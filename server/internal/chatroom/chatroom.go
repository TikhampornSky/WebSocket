package chatroom

import (
	"context"
	"server/server/internal/user"
)

type Chatroom struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Clients []int64 `json:"clients"`
}

type Repository interface {
	CreateChatroom(ctx context.Context, chatroom *Chatroom) (*Chatroom, error)
	JoinChatroom(ctx context.Context, id int64, clientID int64) error
	LeaveChatroom(ctx context.Context, id int64, clientID int64) error
	GetChatroomByID(ctx context.Context, roomId int64) (*GetRoomByID, error)
	UpdateChatroomName(ctx context.Context, id int64, name string) error
	GetAllChatrooms(ctx context.Context) ([]*Chatroom, error)
	DeleteChatroomAll(ctx context.Context) error
}

type GetRoomByID struct {
	ID      int64             `json:"id"`
	Name    string            `json:"name"`
	Clients []user.PublicUser `json:"clients"`
}
