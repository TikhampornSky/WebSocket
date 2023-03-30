# WebSocket

When start project: `make postgres` --> in another commad `docker exec -it postgres15NEW psql` <br> 
To use postgres DB: `make postgres`     `\l`    `\c go-chat`    `\d` (for testing use `go-chat-test`) <br>
To create new migration `migrate create -ext sql -dir db/migrations/ migrationame` <br>
To run server: `go run cmd/main.go` <br>

<br><br>
Create Additional Table Schema <br>

`CREATE TABLE chat_messages ( id SERIAL PRIMARY KEY, sender_id INTEGER NOT NULL REFERENCES users(id), room_id INTEGER NOT NULL REFERENCES chatrooms (id), content TEXT NOT NULL, timestamp TIMESTAMPTZ NOT NULL DEFAULT now() );` <br>

`CREATE TABLE chatrooms ( id bigserial PRIMARY KEY, name varchar NOT NULL UNIQUE );`
`ALTER TABLE chatrooms ADD COLUMN clients BIGINT[] DEFAULT array[]::BIGINT[];`  <br>