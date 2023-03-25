package chatroom

import (
	"context"
	"fmt"
	"log"
	"server/server/db/dbTest"
	"server/server/internal/chatroom"
	"server/server/internal/user"
	"server/server/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var chatroomMockRepo chatroom.Repository
var userMockRepo user.Repository
var dbMock *dbTest.DatabaseTest

func setUpTest() {
	db2, err := dbTest.NewDatabaseTest()
	if err != nil {
		log.Fatalf("Something went wrong. Could not connect to the database. %s", err)
	}
	dbMock = db2
	chatroomMockRepo = chatroom.NewRepository(dbMock.GetDB())
	userMockRepo = user.NewRepository(dbMock.GetDB())
}

func tearDownTest() {
	chatroomMockRepo.DeleteChatroomAll(context.Background())
	userMockRepo.DeleteUserAll(context.Background())
	dbMock.Close()
}

func TestCreateChatroom(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom1",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom.Name, "chatroom1")
}

func TestCreateChatroomDuplicateName(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom2",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom2")

	_, err = chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom2",
	})
	require.ErrorIs(t, err, util.ErrDuplicateChatroom)
}

func TestJoinChatroom(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom3",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom.Name, "chatroom3")

	user, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "joner",
		Email:    "emailJoin",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "joner")
	require.Equal(t, user.Email, "emailJoin")
	require.Equal(t, user.Password, "password")

	err = chatroomMockRepo.JoinChatroom(ctx, chatroom.ID, user.ID)
	require.NoError(t, err)
}

func TestJoinChatroomInvalidUserID(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom3",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom.Name, "chatroom3")

	err = chatroomMockRepo.JoinChatroom(ctx, chatroom.ID, 999999)
	fmt.Println("==> ", err)
	require.ErrorIs(t, err, util.ErrUserIDNotFound)
}
