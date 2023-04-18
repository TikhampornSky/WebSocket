package domain

const (
	Public  string = "public"
	Private        = "private"
)

type Chatroom struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Clients  []int64 `json:"clients"`
	Category string  `json:"category"`
}

type GetRoomByIDRepo struct {
	ID       int64        `json:"id"`
	Name     string       `json:"name"`
	Clients  []PublicUser `json:"clients"`
	Category string       `json:"category"`
}

type CreateChatroomReq struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}
type CreateChatroomRes struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type JoinLeaveChatroomReq struct {
	ID       int64 `json:"id"`
	ClientID int64 `json:"client_id"`
}

type JoinLeaveChatroomRes struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Clients []int64 `json:"clients"`
	Category string `json:"category"`
}

type GetChatroomByIDReq struct {
	ID int64 `json:"id"`
}

type GetChatroomByIDRes struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name"`
	Clients []PublicUser `json:"clients"`
	Category string       `json:"category"`
}

type UpdateChatroomNameReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PublicChatroom struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name"`
	Clients []PublicUser `json:"clients"`
	Category string       `json:"category"`
}
