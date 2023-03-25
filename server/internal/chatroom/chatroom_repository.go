package chatroom

import (
	"context"
	"database/sql"
	"fmt"
	"server/server/util"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &repository{db: db}
}

func (r *repository) CreateChatroom(ctx context.Context, chatroom *Chatroom) (*Chatroom, error) {
	queryFind := "SELECT id FROM chatrooms WHERE name = $1"
	var idFind int64
	err := r.db.QueryRowContext(ctx, queryFind, chatroom.Name).Scan(&idFind)
	if idFind != 0 {
		return &Chatroom{}, util.ErrDuplicateChatroom.With("chatroom with name %s already exists", chatroom.Name)
	}
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("==> ", err)
		return &Chatroom{}, util.ErrInternal.From(err.Error(), err)
	}

	query := "INSERT INTO chatrooms (name) VALUES ($1) RETURNING id"
	var id int64
	err = r.db.QueryRowContext(ctx, query, chatroom.Name).Scan(&id)
	if err != nil {
		return &Chatroom{}, util.ErrInternal.From(err.Error(), err)
	}
	return chatroom, nil
}

func (r *repository) DeleteChatroomAll(ctx context.Context) error { // Testing purposes
	query := "DELETE FROM chatrooms WHERE id > 0"
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return util.ErrInternal.From(err.Error(), err)
	}
	return nil
}

func (r *repository) JoinChatroom(ctx context.Context, id int64, clientID int64) error {
	queryFind := "SELECT id FROM users WHERE id = $1"
	var idFind int64
	err := r.db.QueryRowContext(ctx, queryFind, clientID).Scan(&idFind)
	if idFind == 0 {
		return util.ErrUserIDNotFound.With("user with id %d does not exist", clientID)
	}

	if err != nil {
		return util.ErrInternal.From(err.Error(), err)
	}

	query := "UPDATE chatrooms SET clients = array_append(clients, $1) WHERE id = $2"
	_, err = r.db.ExecContext(ctx, query, clientID, id)
	if err != nil {
		return util.ErrInternal.From(err.Error(), err)
	}
	return nil
}

func (r *repository) GetChatroomByID(ctx context.Context, name string) (*Chatroom, error) {		// Make test
	query := "SELECT id, name, clients FROM chatrooms WHERE name = $1"
	var id int64
	var clients []int64
	err := r.db.QueryRowContext(ctx, query, name).Scan(&id, &name, &clients)
	if err != nil {
		return &Chatroom{}, util.ErrInternal.From(err.Error(), err)
	}
	return &Chatroom{
		ID:      id,
		Name:    name,
		Clients: clients,
	}, nil
}
