package handler

import (
	"fmt"
	"net/http"
	"server/internal/domain"
	"server/internal/port"
	"server/internal/service"
	"server/internal/ws"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub *ws.Hub
	port.ChatroomServicePort
}

func NewWSHandler(hub *ws.Hub, s port.ChatroomServicePort) *WSHandler {
	return &WSHandler{
		hub:                 hub,
		ChatroomServicePort: s,
	}
}

func (h *WSHandler) CreateRoom(c *gin.Context) {
	var req *domain.CreateChatroomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.ChatroomServicePort.CreateChatroom(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[strconv.FormatInt(res.ID, 10)] = &ws.Room{
		ID:      strconv.FormatInt(res.ID, 10),
		Name:    res.Name,
		Clients: make(map[string]*ws.Client),
	}

	c.JSON(http.StatusOK, &domain.Chatroom{
		ID:   res.ID,
		Name: res.Name,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		// origin := r.Header.Get("Origin")
		// return origin == "http://localhost:3000"
	},
}

func (h *WSHandler) JoinRoom(c *gin.Context) {

	// check authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		fmt.Println("unauthorized: no token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	token, err := service.JWTAuthService().ValidateToken(tokenString)
	if token.Valid {
		fmt.Println(token.Claims)
	} else {
		fmt.Println("unauthorized err: ", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// path: /ws/joinRoom/:roomId?userId=123&username=abc
	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clientID, err := strconv.ParseInt(c.Query("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := c.Query("username")

	res, err := h.ChatroomServicePort.JoinChatroom(c.Request.Context(), &domain.JoinLeaveChatroomReq{
		ID:       roomID,
		ClientID: clientID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := &ws.Client{
		Conn:     conn,
		Message:  make(chan *ws.Message),
		ID:       c.Query("userId"),
		RoomID:   c.Param("roomId"),
		Username: username,
	}

	message := &ws.Message{
		Content:  "A new user has joined the room",
		RoomID:   c.Param("roomId"),
		Username: username,
		SenderID: c.Query("userId"),
		Type:     ws.Normal,
	}

	if _, ok := h.hub.Rooms[c.Param("roomId")]; !ok {
		h.hub.Rooms[c.Param("roomId")] = &ws.Room{
			ID:      c.Param("roomId"),
			Name:    res.Name,
			Clients: make(map[string]*ws.Client),
		}
	}

	h.hub.ConnectionMap[client.ID] = client.Conn
	h.hub.BroadcastMap[client.ID] = client.Message

	// Register a new client through the register channel
	h.hub.Register <- client
	// Broadcast the message to all clients in the room
	h.hub.Broadcast <- message

	go client.WriteMessage(h.hub)
	client.ReadMessage(h.hub)
}

func (h *WSHandler) LeaveRoom(c *gin.Context) {

	// path: /ws/leaveRoom/:roomId?userId=123&username=abc
	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clientID, err := strconv.ParseInt(c.Query("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.Query("username")

	err = h.ChatroomServicePort.LeaveChatroom(c.Request.Context(), &domain.JoinLeaveChatroomReq{
		ID:       roomID,
		ClientID: clientID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go ws.LeaveChatroom(h.hub)

	h.hub.LeaveRoom <- &ws.Client{
		Conn:     h.hub.ConnectionMap[c.Query("userId")],
		Message:  h.hub.BroadcastMap[c.Query("userId")],
		ID:       c.Query("userId"),
		RoomID:   c.Param("roomId"),
		Username: username,
	}

	c.JSON(http.StatusOK, nil)
}

func (h *WSHandler) GetRooms(c *gin.Context) {
	rooms := make([]domain.Chatroom, 0)

	arr, err := h.ChatroomServicePort.GetAllChatrooms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, res := range arr {
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rooms = append(rooms, domain.Chatroom{
			ID:      res.ID,
			Name:    res.Name,
			Clients: res.Clients,
		})
	}
	c.JSON(http.StatusOK, rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *WSHandler) GetOnlineClientsInRoom(c *gin.Context) {
	var clients []ClientRes
	roomId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok || h.hub.Rooms[roomId] == nil {
		clients = make([]ClientRes, 0)
	} else {
		for _, c := range h.hub.Rooms[roomId].Clients {
			clients = append(clients, ClientRes{
				ID:       c.ID,
				Username: c.Username,
			})
		}
	}

	c.JSON(http.StatusOK, clients)
}
