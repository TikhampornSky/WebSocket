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
