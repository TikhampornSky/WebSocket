package main

import (
	"log"
	"server/server/db"
	user "server/server/internal"
	"server/server/internal/ws"
	"server/server/router"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Something went wrong. Could not connect to the database. %s", err)
	}

	userRepo := user.NewRepository(db.GetDB())
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")

	defer db.Close()
}
