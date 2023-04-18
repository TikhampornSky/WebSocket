# WebSocket

# ==============================================
To start project: 
1. Delete postgresql folder
2. `make postgresinit`
3. `make createdb` (to create database)
4. `make createdbtest` (to create database for testing)
5. `make migrateup` (create tables in database)
6. `make migrateuptest` (create tables in database testing)
7. `make postgres` (to start project) <-- 
# ==============================================

When start project: `make postgres` --> in another commad `docker exec -it postgres15NEW psql` <br> 
To use postgres DB: `make postgres`     `\l`    `\c go-chat`    `\d` (for testing use `go-chat-test`) <br>
To create new migration `migrate create -ext sql -dir db/migrations/ migrationame` <br>
To run server: `go run cmd/main.go` <br>

<br><br>
Create Additional Table Schema <br>

`CREATE TABLE chat_messages ( id SERIAL PRIMARY KEY, sender_id INTEGER NOT NULL REFERENCES users(id), room_id INTEGER NOT NULL REFERENCES chatrooms (id), content TEXT NOT NULL, timestamp TIMESTAMPTZ NOT NULL DEFAULT now() );` <br>

`CREATE TABLE chatrooms ( id bigserial PRIMARY KEY, name varchar NOT NULL UNIQUE );` <br>
`ALTER TABLE chatrooms ADD COLUMN clients BIGINT[] DEFAULT array[]::BIGINT[];`  <br>

`migrate create -ext sql -dir db/migrations name_of_migrate` <br>
`\dT+ roomType` <br>