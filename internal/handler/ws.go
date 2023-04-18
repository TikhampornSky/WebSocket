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
	"github.com/golang-jwt/jwt/v4"
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
	if req.Category != domain.Public && req.Category != domain.Private {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category must be public or private"})
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

	c.JSON(http.StatusCreated, &domain.Chatroom{
		ID:   res.ID,
		Name: res.Name,
		Category: res.Category,
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WSHandler) JoinRoom(c *gin.Context) {

	// check authorization
	tokenString := c.GetHeader("Sec-Websocket-Protocol")
	fmt.Println("tokenString: ", tokenString)
	if tokenString == "" {
		fmt.Println("unauthorized: no token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	token, err := service.JWTAuthService().ValidateToken(tokenString)
	if token.Valid {
		c.Set("userID", token.Claims.(jwt.MapClaims)["id"])
		c.Set("username", token.Claims.(jwt.MapClaims)["username"])
	} else {
		fmt.Println("unauthorized err: ", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, http.Header{
		"Sec-websocket-Protocol": websocket.Subprotocols(c.Request),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	clientID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := c.MustGet("username").(string)

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
		ID:       userID,
		RoomID:   c.Param("roomId"),
		Username: username,
	}

	message := &ws.Message{
		Content:  "A new user has joined the room",
		RoomID:   c.Param("roomId"),
		Username: username,
		SenderID: userID,
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

	roomID, err := strconv.ParseInt(c.Param("roomId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("userID").(string)
	clientID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := c.MustGet("username").(string)

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
		Conn:     h.hub.ConnectionMap[userID],
		Message:  h.hub.BroadcastMap[userID],
		ID:       userID,
		RoomID:   c.Param("roomId"),
		Username: username,
	}

	c.JSON(http.StatusOK, nil)
}

func (h *WSHandler) GetRooms(c *gin.Context) {
	rooms := make([]domain.Chatroom, 0)

	userID := c.MustGet("userID").(string)
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	arr, err := h.ChatroomServicePort.GetAllChatrooms(c.Request.Context(), userIDInt)
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
			Category: res.Category,
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
