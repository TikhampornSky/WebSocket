package chatroom

import (
	"context"
	"server/server/internal/user"
)

type Chatroom struct {
	ID      uint8   `json:"id"`
	Name    string  `json:"name"`
	Clients []uint8 `json:"clients"`
}

type Repository interface {
	CreateChatroom(ctx context.Context, chatroom *Chatroom) (*Chatroom, error)
	JoinChatroom(ctx context.Context, id uint8, clientID int64) error
	GetChatroomByID(ctx context.Context, roomId uint8) (*GetRoomByID, error)
	DeleteChatroomAll(ctx context.Context) error
	// UpdateChatroomName(ctx context.Context, id uint8, name string) error
	// GetAllChatrooms(ctx context.Context) ([]*Chatroom, error)
}

type GetRoomByID struct {
	ID      uint8             `json:"id"`
	Name    string            `json:"name"`
	Clients []user.PublicUser `json:"clients"`
}
