package domain

type Chatroom struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Clients []int64 `json:"clients"`
}

type GetRoomByIDRepo struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name"`
	Clients []PublicUser `json:"clients"`
}

type CreateChatroomReq struct {
	Name string `json:"name"`
}
type CreateChatroomRes struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type JoinLeaveChatroomReq struct {
	ID       int64 `json:"id"`
	ClientID int64 `json:"client_id"`
}

type JoinLeaveChatroomRes struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Clients []int64 `json:"clients"`
}

type GetChatroomByIDReq struct {
	ID int64 `json:"id"`
}

type GetChatroomByIDRes struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name"`
	Clients []PublicUser `json:"clients"`
}

type UpdateChatroomNameReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PublicChatroom struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name"`
	Clients []PublicUser `json:"clients"`
}
