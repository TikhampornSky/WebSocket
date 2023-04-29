FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod verify

COPY . .

EXPOSE 8080

ENV POSTGRES_USER=root
ENV POSTGRES_PASSWORD=password
ENV POSTGRES_DB=go-chat
ENV POSTGRES_HOST=localhost

RUN go build -o main cmd/main.go

CMD ["./main"]
