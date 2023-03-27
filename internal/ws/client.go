package ws

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type MessageType int

const (
    Normal MessageType = iota
    LeaveRoom
)

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
	SenderID string `json:"senderId"`
	Type     MessageType `json:"type"`
}

func (c *Client) WriteMessage(h *Hub) {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <- c.Message
		if message == nil {		// When Leaving room, the message is nil
			return
		}
		if !ok {
			fmt.Println("Write error", ok)
			return
		}

		c.Conn.WriteJSON(message)

		if message.Type == LeaveRoom && message.SenderID == c.ID {		// When Leaving room, close the channel and delete the client from the room
			close(h.BroadcastMap[message.SenderID])		
			delete(h.BroadcastMap, message.SenderID)
			delete(h.Rooms[message.RoomID].Clients, message.SenderID)
			delete(h.ConnectionMap, message.SenderID)
			log.Println("Deleted client", message.SenderID, "from room", message.RoomID)
		}
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
			SenderID: c.ID,
			Type:    Normal,
		}
		hub.Broadcast <- msg
	}
}

func LeaveChatroom(hub *Hub) {
	for {
		client, ok := <-hub.LeaveRoom
		if !ok {
			fmt.Println("Leave error", ok)
			return
		}
		hub.Unregister <- client
	}
}
