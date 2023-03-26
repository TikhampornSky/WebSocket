# WebSocket

When start project: `cd server` `make postgres` --> in another commad `docker exec -it postgres15 psql` <br> 
To use postgres DB: `make postgres`     `\l`    `\c go-chat`    `\d` <br>
To create new migration `cd server` `migrate create -ext sql -dir db/migrations/ migrationame` <br>
To run server: `cd server`  `go run cmd/main.go` <br>

<br><br>
Create Additional Table Schema <br>

`CREATE TABLE chat_messages ( id SERIAL PRIMARY KEY, sender_id INTEGER NOT NULL REFERENCES users(id), room_id INTEGER NOT NULL REFERENCES chatroom (id), content TEXT NOT NULL, timestamp TIMESTAMPTZ NOT NULL DEFAULT now() );` <br>

`CREATE TABLE chatrooms ( id bigserial PRIMARY KEY, name varchar NOT NULL UNIQUE );`
`ALTER TABLE chatrooms ADD COLUMN clients BIGINT[] DEFAULT array[]::BIGINT[];`  <br>