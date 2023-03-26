package repo_test

import (
	"context"
	"server/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateChatroom(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom1",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom.Name, "chatroom1")

	res, err := chatroomMockRepo.GetChatroomByID(context.Background(), chatroom.ID)
	require.NoError(t, err)
	require.Equal(t, res.Name, "chatroom1")
}

func TestCreateChatroomDuplicateName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom2",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom2")

	_, err = chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom2",
	})
	require.ErrorIs(t, err, domain.ErrDuplicateChatroom)
}
func TestJoinChatroom(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom3",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom3")

	chatroom2, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "chatroom3")

	user, err := userMockRepo.CreateUser(ctx, &domain.User{
		Username: "joner",
		Email:    "emailJoin",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "joner")
	require.Equal(t, user.Email, "emailJoin")
	require.Equal(t, user.Password, "password")

	_, err = chatroomMockRepo.JoinChatroom(ctx, chatroom1.ID, user.ID)
	require.NoError(t, err)
}

func TestJoinChatroomInvalidUserID(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom4",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom.Name, "chatroom4")

	_, err = chatroomMockRepo.JoinChatroom(ctx, chatroom.ID, 0)
	require.ErrorIs(t, err, domain.ErrUserIDNotFound)
}

func TestJoinChatroomInvalidChatroomID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := userMockRepo.CreateUser(ctx, &domain.User{
		Username: "joner1",
		Email:    "emailJoin1",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "joner1")
	require.Equal(t, user.Email, "emailJoin1")
	require.Equal(t, user.Password, "password")

	_, err = chatroomMockRepo.JoinChatroom(ctx, 0, user.ID)
	require.ErrorIs(t, err, domain.ErrChatroomIDNotFound)
}

func TestGetChatroomByID(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom5",
	})
	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom5")

	user1, err := userMockRepo.CreateUser(ctx, &domain.User{
		Username: "jonerr",
		Email:    "emailJoinn",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user1.Username, "jonerr")
	require.Equal(t, user1.Email, "emailJoinn")
	require.Equal(t, user1.Password, "password")
	user2, err := userMockRepo.CreateUser(ctx, &domain.User{
		Username: "joner2",
		Email:    "emailJoin2",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user2.Username, "joner2")
	require.Equal(t, user2.Email, "emailJoin2")
	require.Equal(t, user2.Password, "password")

	_, err = chatroomMockRepo.JoinChatroom(ctx, chatroom1.ID, user1.ID)
	require.NoError(t, err)
	_, err = chatroomMockRepo.JoinChatroom(ctx, chatroom1.ID, user2.ID)
	require.NoError(t, err)

	chatroom2, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "chatroom5")
	for _, user := range chatroom2.Clients {
		if user.ID == user1.ID {
			require.Equal(t, user.Username, "jonerr")
			require.Equal(t, user.Email, "emailJoinn")
		} else {
			require.Equal(t, user.Username, "joner2")
			require.Equal(t, user.Email, "emailJoin2")
		}
	}
}

func TestGetChatroomByIDNoClients(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom6",
	})
	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom6")

	chatroom2, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "chatroom6")
}

func TestGetChatroomByIDInvalidID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := chatroomMockRepo.GetChatroomByID(ctx, 0)
	require.ErrorIs(t, err, domain.ErrChatroomIDNotFound)
}

func TestUpdateChatroomName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := chatroomMockRepo.UpdateChatroomName(ctx, 0, "newChatRoomName")
	require.ErrorIs(t, err, domain.ErrChatroomIDNotFound)
}

func TestGetAllChatrooms(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroomMockRepo.DeleteChatroomAll(ctx)

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom8",
	})
	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom8")

	chatroom2, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroomMockRepo.DeleteChatroomAll(ctx)

	chatrooms, err := chatroomMockRepo.GetAllChatrooms(ctx)
	require.NoError(t, err)
	require.Equal(t, len(chatrooms), 0)
}

func TestLeaveChatroom(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroomMockRepo.DeleteChatroomAll(ctx)

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom3",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom3")

	chatroom2, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom2.Name, "chatroom3")

	user, err := userMockRepo.CreateUser(ctx, &domain.User{
		Username: "joner3",
		Email:    "emailJoin3",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "joner3")
	require.Equal(t, user.Email, "emailJoin3")
	require.Equal(t, user.Password, "password")

	_, err = chatroomMockRepo.JoinChatroom(ctx, chatroom1.ID, user.ID)
	require.NoError(t, err)

	err = chatroomMockRepo.LeaveChatroom(ctx, chatroom1.ID, user.ID)
	require.NoError(t, err)

	chatroom3, err := chatroomMockRepo.GetChatroomByID(ctx, chatroom1.ID)
	require.NoError(t, err)
	require.Equal(t, chatroom3.Name, "chatroom3")
	require.Equal(t, len(chatroom3.Clients), 0)
}

func TestLeaveChatroomInvalidChatroomID(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := userMockRepo.CreateUser(ctx, &domain.User{
		Username: "joner4",
		Email:    "emailJoin4",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "joner4")
	require.Equal(t, user.Email, "emailJoin4")
	require.Equal(t, user.Password, "password")

	err = chatroomMockRepo.LeaveChatroom(ctx, 0, user.ID)
	require.ErrorIs(t, err, domain.ErrChatroomIDNotFound)
}

func TestLeaveChatroomInvalidUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatroomMockRepo.DeleteChatroomAll(ctx)

	chatroom1, err := chatroomMockRepo.CreateChatroom(ctx, &domain.Chatroom{
		Name: "chatroom4",
	})

	require.NoError(t, err)
	require.Equal(t, chatroom1.Name, "chatroom4")

	err = chatroomMockRepo.LeaveChatroom(ctx, chatroom1.ID, 0)
	require.ErrorIs(t, err, domain.ErrUserIDNotFound)
}
