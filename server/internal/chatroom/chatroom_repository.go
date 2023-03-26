package chatroom

import (
	"context"
	"database/sql"
	"server/server/internal/user"
	"server/server/util"
)

type DBTXChat interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTXChat
}

func NewRepository(db DBTXChat) Repository {
	return &repository{db: db}
}

func (r *repository) CreateChatroom(ctx context.Context, chatroom *Chatroom) (*Chatroom, error) {
	queryFind := "SELECT id FROM chatrooms WHERE name = $1"
	var idFind uint8
	err := r.db.QueryRowContext(ctx, queryFind, chatroom.Name).Scan(&idFind)
	if err != sql.ErrNoRows {
		return &Chatroom{}, util.ErrDuplicateChatroom.With("chatroom with name %s already exists", chatroom.Name)
	}
	if err != nil && err != sql.ErrNoRows {
		return &Chatroom{}, util.ErrInternal.From(err.Error(), err)
	}

	query := "INSERT INTO chatrooms (name) VALUES ($1) RETURNING id"
	var id uint8
	err = r.db.QueryRowContext(ctx, query, chatroom.Name).Scan(&id)
	if err != nil {
		return &Chatroom{}, util.ErrInternal.From(err.Error(), err)
	}
	chatroom.ID = id
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

func (r *repository) JoinChatroom(ctx context.Context, id uint8, clientID int64) error {
	queryFindUser := "SELECT id FROM users WHERE id = $1"
	var idFindUser int64
	err := r.db.QueryRowContext(ctx, queryFindUser, clientID).Scan(&idFindUser)
	if err == sql.ErrNoRows {
		return util.ErrUserIDNotFound.With("user with id %d does not exist", clientID)
	}
	if err != nil {
		return util.ErrInternal.From(err.Error(), err)
	}

	var resId uint8
	query := "UPDATE chatrooms SET clients = array_append(clients, $1) WHERE id = $2 AND NOT ($1 = ANY(clients)) RETURNING id"
	err = r.db.QueryRowContext(ctx, query, clientID, id).Scan(&resId)
	if err == sql.ErrNoRows {
		return util.ErrChatroomIDNotFound.With("chatroom with id %d does not exist", id)
	}
	if err != nil {
		return util.ErrInternal.From(err.Error(), err)
	}
	return nil
}

func (r *repository) GetChatroomByID(ctx context.Context, roomId uint8) (*GetRoomByID, error) {
	query := `SELECT chatrooms.id, name as roomName, clients, users.id as userId, username, email
				FROM chatrooms LEFT JOIN users ON users.id = ANY (chatrooms.clients) 
				WHERE chatrooms.id = $1 
				ORDER BY users.id;`
	rows, err := r.db.QueryContext(ctx, query, roomId)
	if err != nil {
		return &GetRoomByID{}, util.ErrInternal.From(err.Error(), err)
	}

	var chatroom GetRoomByID
	var clients []user.PublicUser
	for rows.Next() {
		var userid sql.NullInt64
		var username sql.NullString
		var email sql.NullString
		var chatroomTmp Chatroom

		err = rows.Scan(&chatroomTmp.ID, &chatroomTmp.Name, &chatroomTmp.Clients, &userid, &username, &email)
		
		if userid.Valid {
			chatroom.ID = chatroomTmp.ID
			chatroom.Name = chatroomTmp.Name
			clients = append(clients, user.PublicUser{
				ID:       userid.Int64,
				Username: username.String,
				Email:    email.String,
			})
		} else {
			chatroom.ID = chatroomTmp.ID
			chatroom.Name = chatroomTmp.Name
		}
		if err != nil {
			return &GetRoomByID{}, util.ErrInternal.From(err.Error(), err)
		}
	}

	chatroom.Clients = clients

	if chatroom.ID == 0 {
		return &GetRoomByID{}, util.ErrChatroomIDNotFound.With("chatroom with id %d does not exist", roomId)
	}
	if err != nil {
		return &GetRoomByID{}, util.ErrInternal.From(err.Error(), err)
	}

	return &chatroom, nil
}