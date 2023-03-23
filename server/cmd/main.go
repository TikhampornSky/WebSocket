package main

import (
	"log"
	"server/server/db"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Something went wrong. Could not connect to the database. %s", err)
	}
	defer db.Close()
}
