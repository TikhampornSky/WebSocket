package repo

import (
	"context"
	"database/sql"
	"server/internal/domain"
	"server/internal/port"

	"github.com/lib/pq"
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

func NewChatroomRepository(db DBTXChat) port.ChatroomRepoPort {
	return &repository{db: db}
}

func (r *repository) CreateChatroom(ctx context.Context, chatroom *domain.Chatroom) (*domain.Chatroom, error) {
	queryFind := "SELECT id FROM chatrooms WHERE name = $1"
	var idFind int64
	err := r.db.QueryRowContext(ctx, queryFind, chatroom.Name).Scan(&idFind)
	if err != sql.ErrNoRows {
		return &domain.Chatroom{}, domain.ErrDuplicateChatroom.With("chatroom with name %s already exists", chatroom.Name)
	}
	if err != nil && err != sql.ErrNoRows {
		return &domain.Chatroom{}, domain.ErrInternal.From(err.Error(), err)
	}

	query := "INSERT INTO chatrooms (name) VALUES ($1) RETURNING id"
	var id int64
	err = r.db.QueryRowContext(ctx, query, chatroom.Name).Scan(&id)
	if err != nil {
		return &domain.Chatroom{}, domain.ErrInternal.From(err.Error(), err)
	}
	chatroom.ID = id
	return chatroom, nil
}

func (r *repository) DeleteChatroomAll(ctx context.Context) error { // Testing purposes
	query := "DELETE FROM chatrooms WHERE id > 0"
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}
	return nil
}

func (r *repository) JoinChatroom(ctx context.Context, id int64, clientID int64) error {
	queryFindUser := "SELECT id FROM users WHERE id = $1"
	var idFindUser int64
	err := r.db.QueryRowContext(ctx, queryFindUser, clientID).Scan(&idFindUser)
	if err == sql.ErrNoRows {
		return domain.ErrUserIDNotFound.With("user with id %d does not exist", clientID)
	}
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}

	var resId int64
	query := "UPDATE chatrooms SET clients = array_append(clients, $1) WHERE id = $2 AND NOT ($1 = ANY(clients)) RETURNING id"
	err = r.db.QueryRowContext(ctx, query, clientID, id).Scan(&resId)
	if err == sql.ErrNoRows {
		return domain.ErrChatroomIDNotFound.With("chatroom with id %d does not exist", id)
	}
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}
	return nil
}

func (r *repository) LeaveChatroom(ctx context.Context, id int64, clientID int64) error {
	queryFindUser := "SELECT id FROM users WHERE id = $1"
	var idFindUser int64
	err := r.db.QueryRowContext(ctx, queryFindUser, clientID).Scan(&idFindUser)
	if err == sql.ErrNoRows {
		return domain.ErrUserIDNotFound.With("user with id %d does not exist", clientID)
	}
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}

	var resId int64
	query := "UPDATE chatrooms SET clients = array_remove(clients, $1) WHERE id = $2 AND ($1 = ANY(clients)) RETURNING id"
	err = r.db.QueryRowContext(ctx, query, clientID, id).Scan(&resId)
	if err == sql.ErrNoRows {
		return domain.ErrChatroomIDNotFound.With("chatroom with id %d does not exist", id)
	}
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}
	return nil
}

func (r *repository) GetChatroomByID(ctx context.Context, roomId int64) (*domain.GetRoomByIDRepo, error) {
	query := `SELECT chatrooms.id, name as roomName, clients, users.id as userId, username, email
				FROM chatrooms LEFT JOIN users ON users.id = ANY (chatrooms.clients) 
				WHERE chatrooms.id = $1 
				ORDER BY users.id;`
	rows, err := r.db.QueryContext(ctx, query, roomId)
	if err != nil {
		return &domain.GetRoomByIDRepo{}, domain.ErrInternal.From(err.Error(), err)
	}

	var chatroomByID domain.GetRoomByIDRepo
	var clients []domain.PublicUser
	for rows.Next() {
		var userid sql.NullInt64
		var username sql.NullString
		var email sql.NullString
		var chatroomTmp domain.Chatroom
		err = rows.Scan(&chatroomTmp.ID, &chatroomTmp.Name, pq.Array(&chatroomTmp.Clients), &userid, &username, &email)

		if userid.Valid {
			chatroomByID.ID = chatroomTmp.ID
			chatroomByID.Name = chatroomTmp.Name
			clients = append(clients, domain.PublicUser{
				ID:       userid.Int64,
				Username: username.String,
				Email:    email.String,
			})
		} else {
			chatroomByID.ID = chatroomTmp.ID
			chatroomByID.Name = chatroomTmp.Name
		}

		if err != nil {
			return &domain.GetRoomByIDRepo{}, domain.ErrInternal.From(err.Error(), err)
		}
	}

	chatroomByID.Clients = clients

	if chatroomByID.ID == 0 {
		return &domain.GetRoomByIDRepo{}, domain.ErrChatroomIDNotFound.With("chatroom with id %d does not exist", roomId)
	}
	if err != nil {
		return &domain.GetRoomByIDRepo{}, domain.ErrInternal.From(err.Error(), err)
	}

	return &chatroomByID, nil
}

func (r *repository) UpdateChatroomName(ctx context.Context, id int64, name string) error {
	query := "UPDATE chatrooms SET name = $1 WHERE id = $2 RETURNING id"
	var resId int64
	err := r.db.QueryRowContext(ctx, query, name, id).Scan(&resId)
	if err == sql.ErrNoRows {
		return domain.ErrChatroomIDNotFound.With("chatroom with id %d does not exist", id)
	}
	if err != nil {
		return domain.ErrInternal.From(err.Error(), err)
	}
	return nil
}

func (r *repository) GetAllChatrooms(ctx context.Context) ([]*domain.Chatroom, error) {
	query := "SELECT id, name, clients FROM chatrooms"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return []*domain.Chatroom{}, domain.ErrInternal.From(err.Error(), err)
	}

	var chatrooms []*domain.Chatroom
	for rows.Next() {
		var chatroom domain.Chatroom
		err = rows.Scan(&chatroom.ID, &chatroom.Name, pq.Array(&chatroom.Clients))
		if err != nil {
			return []*domain.Chatroom{}, domain.ErrInternal.From(err.Error(), err)
		}

		chatrooms = append(chatrooms, &chatroom)
	}
	return chatrooms, nil
}
