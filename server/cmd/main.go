package main

import (
	"log"
	"server/server/db"
	user "server/server/internal"
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

	router.InitRouter(userHandler)
	router.Start("0.0.0.0:8080")

	defer db.Close()
}
