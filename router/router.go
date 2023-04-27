package router

import (
	"server/internal/handler"
	"server/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(userHandler *handler.UserHandler, wsHandler *handler.WSHandler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "PUT"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Origin", "Accept", "X-Requested-With", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Allow-Credentials"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	r.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)

	r.Use(middleware.AuthorizeJWT())
	{
		r.GET("/users", userHandler.GetAllUsers)
		r.PATCH("/user/self", userHandler.UpdateUsername)
		r.PATCH("/user/self/password", userHandler.UpdatePassword)
		r.POST("/ws/createRoom", wsHandler.CreateRoom)
		r.POST("/ws/createDM", wsHandler.CreateDM)
		r.GET("/ws/leaveRoom/:roomId", wsHandler.LeaveRoom)
		r.GET("/ws/getRooms", wsHandler.GetRooms)
		r.GET("/ws/getDMs", wsHandler.GetDMs)
		r.GET("/ws/getClients/:roomId", wsHandler.GetOnlineClientsInRoom) // Only show client that are now online (join the room) in the new connection
	}
}

func Start(addr string) error {
	return r.Run(addr)
}
