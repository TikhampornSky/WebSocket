package main

import (
	"log"
	"server/db"
	"server/internal/handler"
	"server/internal/repo"
	"server/internal/service"
	"server/internal/ws"
	"server/router"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Something went wrong. Could not connect to the database. %s", err)
	}

	userRepo := repo.NewUserRepository(db.GetDB())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")

	defer db.Close()
}
