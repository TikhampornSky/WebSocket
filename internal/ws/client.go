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

type Message struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
	SenderID string `json:"senderId"`
}

func (c *Client) WriteMessage(h *Hub) {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <- c.Message
		if !ok {
			fmt.Println("Write error", ok)
			fmt.Println("Messade", message) 		// nil (Why! )
			return
		}

		log.Println("Write message", message.Content, " --from: ", message.SenderID, " -Receriver ID-> ", c.ID)
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
		}
		log.Println("Read message", msg)
		hub.Broadcast <- msg
	}
}

func LeaveChatroom(hub *Hub) {
	for {
		client, ok := <-hub.LeaveRoom
		// log.Panicln("Leave chatroom--> ", client, ok)

		// defer hub.Connection[client.ID].Close()

		if !ok {
			fmt.Println("Leave error", ok)
			return
		}

		hub.Unregister <- client
	}
}
