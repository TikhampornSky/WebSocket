# WebSocket

When start project: `cd server` `make postgres` --> in another commad `docker exec -it postgres15 psql` <br> 
To use postgres DB: `make postgres`     `\l`    `\c go-chat`    `\d` <br>
To create new migration `cd server` `migrate create -ext sql -dir db/migrations/ migrationame` <br>
To run server: `cd server`  `go run cmd/main.go` <br>