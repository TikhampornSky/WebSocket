package ws

import (
	"github.com/gorilla/websocket"
)

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms         map[string]*Room
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan *Message
	LeaveRoom     chan *Client
	ConnectionMap map[string]*websocket.Conn
	BroadcastMap  map[string]chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:         make(map[string]*Room),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan *Message, 5),
		LeaveRoom:     make(chan *Client),
		ConnectionMap: make(map[string]*websocket.Conn),
		BroadcastMap:  make(map[string]chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.RoomID]; ok {
				room := h.Rooms[client.RoomID]

				if _, ok := room.Clients[client.ID]; !ok { // if client is not in the room
					room.Clients[client.ID] = client
				}
			}
		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.RoomID]; ok {
				if _, ok := h.Rooms[client.RoomID].Clients[client.ID]; ok {
					if len(h.Rooms[client.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{ // Broadcast a message saying that the user has left the room
							Content:  client.Username + " left the room",
							RoomID:   client.RoomID,
							Username: client.Username,
							SenderID: client.ID,
							Type:     LeaveRoom,
						}
					}
				}
			}

		case message := <-h.Broadcast:
			if _, ok := h.Rooms[message.RoomID]; ok {
				for _, cl := range h.Rooms[message.RoomID].Clients {
					cl.Message <- message
				}
			}
		}
	}
}
