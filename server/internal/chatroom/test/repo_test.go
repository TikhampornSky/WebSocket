package chatroom

import (
	"context"
	"log"
	"server/server/dbTest"
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

	res, err := chatroomMockRepo.GetChatroomByID(context.Background(), chatroom.ID)
	require.NoError(t, err)
	require.Equal(t, res.Name, "chatroom1")
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

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom3",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom3")

	chatroom2, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "chatroom3")

	user, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "joner",
		Email:    "emailJoin",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "joner")
	require.Equal(t, user.Email, "emailJoin")
	require.Equal(t, user.Password, "password")

	err = chatroomMockRepo.JoinChatroom(ctx, chatroom1.ID, user.ID)
	require.NoError(t, err)
}

func TestJoinChatroomInvalidUserID(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom4",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom.Name, "chatroom4")

	err = chatroomMockRepo.JoinChatroom(ctx, chatroom.ID, 0)
	require.ErrorIs(t, err, util.ErrUserIDNotFound)
}

func TestJoinChatroomInvalidChatroomID(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "joner1",
		Email:    "emailJoin1",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "joner1")
	require.Equal(t, user.Email, "emailJoin1")
	require.Equal(t, user.Password, "password")

	err = chatroomMockRepo.JoinChatroom(ctx, 0, user.ID)
	require.ErrorIs(t, err, util.ErrChatroomIDNotFound)
}

func TestGetChatroomByID(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom5",
	})
	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom5")

	user1, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "joner",
		Email:    "emailJoin",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user1.Username, "joner")
	require.Equal(t, user1.Email, "emailJoin")
	require.Equal(t, user1.Password, "password")
	user2, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "joner2",
		Email:    "emailJoin2",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user2.Username, "joner2")
	require.Equal(t, user2.Email, "emailJoin2")
	require.Equal(t, user2.Password, "password")

	err = chatroomMockRepo.JoinChatroom(ctx, chatroom1.ID, user1.ID)
	require.NoError(t, err)
	err = chatroomMockRepo.JoinChatroom(ctx, chatroom1.ID, user2.ID)
	require.NoError(t, err)

	chatroom2, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "chatroom5")
	for _, user := range chatroom2.Clients {
		if user.ID == user1.ID {
			require.Equal(t, user.Username, "joner")
			require.Equal(t, user.Email, "emailJoin")
		} else {
			require.Equal(t, user.Username, "joner2")
			require.Equal(t, user.Email, "emailJoin2")
		}
	}
}

func TestGetChatroomByIDNoClients(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom6",
	})
	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom6")

	chatroom2, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "chatroom6")
}

func TestGetChatroomByIDInvalidID(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := chatroomMockRepo.GetChatroomByID(ctx, 0)
	require.ErrorIs(t, err, util.ErrChatroomIDNotFound)
}

func TestUpdateChatroomName(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom7",
	})
	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom7")

	err = chatroomMockRepo.UpdateChatroomName(ctx, chatroom1.ID, "newChatRoomName")
	require.NoError(t, err)

	chatroom2, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "newChatRoomName")
}

func TestUpdateChatroomNameInvalidRoomId(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := chatroomMockRepo.UpdateChatroomName(ctx, 0, "newChatRoomName")
	require.ErrorIs(t, err, util.ErrChatroomIDNotFound)
}

func TestGetAllChatrooms(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom8",
	})
	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom8")

	chatroom2, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom9",
	})
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "chatroom9")

	chatrooms, err := chatroomMockRepo.GetAllChatrooms(ctx)
	require.NoError(t, err)
	require.Equal(t, len(chatrooms), 2)
	for _, chatroom := range chatrooms {
		if chatroom.ID == chatroom1.ID {
			require.Equal(t, chatroom.Name, "chatroom8")
		} else {
			require.Equal(t, chatroom.Name, "chatroom9")
		}
	}
}

func TestGetAllChatroomsNoChatrooms(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatrooms, err := chatroomMockRepo.GetAllChatrooms(ctx)
	require.NoError(t, err)
	require.Equal(t, len(chatrooms), 0)
}

func TestLeaveChatroom(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom3",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom3")

	chatroom2, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "chatroom3")

	user, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "joner3",
		Email:    "emailJoin3",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "joner3")
	require.Equal(t, user.Email, "emailJoin")
	require.Equal(t, user.Password, "password")

	err = chatroomMockRepo.JoinChatroom(ctx, chatroom1.ID, user.ID)
	require.NoError(t, err)

	err = chatroomMockRepo.LeaveChatroom(ctx, chatroom1.ID, user.ID)
	require.NoError(t, err)

	chatroom3, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom3.Name, "chatroom3")
	require.Equal(t, len(chatroom3.Clients), 0)
}

func TestLeaveChatroomInvalidChatroomID(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "joner4",
		Email:    "emailJoin4",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "joner4")
	require.Equal(t, user.Email, "emailJoin4")
	require.Equal(t, user.Password, "password")

	err = chatroomMockRepo.LeaveChatroom(ctx, 0, user.ID)
	require.ErrorIs(t, err, util.ErrChatroomIDNotFound)
}

func TestLeaveChatroomInvalidUserID(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &chatroom.Chatroom{
		Name: "chatroom4",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom4")

	err = chatroomMockRepo.LeaveChatroom(ctx, chatroom1.ID, 0)
	require.ErrorIs(t, err, util.ErrUserIDNotFound)
}