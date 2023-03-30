postgresinit:
	docker run --name postgres15NEW -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d -v /Users/tikhamporntepsut/Documents/go_socket/postgresql:/var/lib/postgresql/data postgres:15-alpine

postgres:
	docker exec -it postgres15NEW psql

createdb:
	docker exec -it postgres15NEW createdb --username=root --owner=root go-chat

dropdb:
	docker exec -it postgres15NEW dropdb go-chat

migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose down

.PHONY: postgresinit postgres createdb dropdb