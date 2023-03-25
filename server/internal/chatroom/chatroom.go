package chatroom

import "context"

type Chatroom struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Clients []int64 `json:"clients"`
}

type Repository interface {
	CreateChatroom(ctx context.Context, chatroom *Chatroom) (*Chatroom, error)
	JoinChatroom(ctx context.Context, id int64, clientID int64) error
	GetChatroomByID(ctx context.Context, name string) (*Chatroom, error)
	DeleteChatroomAll(ctx context.Context) error
	// UpdateChatroomName(ctx context.Context, id int64, name string) error
	// GetAllChatrooms(ctx context.Context) ([]*Chatroom, error)
}
