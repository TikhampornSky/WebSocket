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
	LeaveChatroom
)

type Message struct {
	Content  string      `json:"content"`
	RoomID   string      `json:"roomId"`
	Username string      `json:"username"`
	SenderID string      `json:"senderId"`
	Type     MessageType `json:"type"`
}

func (c *Client) WriteMessage(h *Hub) {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			fmt.Println("Write error", ok)
			return
		}

		log.Println("Write message", message, " -Receriver ID-> ", c.ID)
		c.Conn.WriteJSON(message)
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
			Type:     Normal,
		}
		log.Println("Read message", msg)
		hub.Broadcast <- msg
	}
}
