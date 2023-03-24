package user

import (
	"context"
	"log"
	"server/server/db/dbTest"
	"server/server/internal/user"
	"server/server/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var userMockRepo user.Repository
var dbMock *dbTest.DatabaseTest

func setUpTest() {
	db2, err := dbTest.NewDatabaseTest()
	if err != nil {
		log.Fatalf("Something went wrong. Could not connect to the database. %s", err)
	}
	dbMock = db2
	userMockRepo = user.NewRepository(dbMock.GetDB())
}

func tearDownTest() {
	userMockRepo.DeleteUserAll(context.Background())
	dbMock.Close()
}


func TestCreateUser(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "username",
		Email:    "email",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "username")
	require.Equal(t, user.Email, "email")
	require.Equal(t, user.Password, "password")
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user1, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "username",
		Email:    "email",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user1.Username, "username")
	require.Equal(t, user1.Email, "email")
	require.Equal(t, user1.Password, "password")

	_, err = userMockRepo.CreateUser(ctx, &user.User{
		Username: "username1",
		Email:    "email",
		Password: "password",
	})

	require.ErrorIs(t, err, util.ErrDuplicateEmail)
}

func TestCreateUserDuplicateUsername(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user1, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "username",
		Email:    "email",
		Password: "password",
	})

	require.NoError(t, err)
	require.Equal(t, user1.Username, "username")
	require.Equal(t, user1.Email, "email")
	require.Equal(t, user1.Password, "password")

	_, err = userMockRepo.CreateUser(ctx, &user.User{
		Username: "username",
		Email:    "email1",
		Password: "password",
	})

	require.ErrorIs(t, err, util.ErrDuplicateUsername)
}

func TestGetUserByEmail(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "username2",
		Email:    "email2",
		Password: "password2",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "username2")
	require.Equal(t, user.Email, "email2")
	require.Equal(t, user.Password, "password2")

	user, err = userMockRepo.GetUserByEmail(ctx, "email2")
	require.NoError(t, err)
	require.Equal(t, user.Username, "username2")
	require.Equal(t, user.Email, "email2")
	require.Equal(t, user.Password, "password2")
}

func TestUpdateUsername(t *testing.T) {
	setUpTest()
	defer tearDownTest()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := userMockRepo.CreateUser(ctx, &user.User{
		Username: "username3",
		Email:    "email3",
		Password: "password3",
	})

	require.NoError(t, err)
	require.Equal(t, user.Username, "username3")
	require.Equal(t, user.Email, "email3")
	require.Equal(t, user.Password, "password3")

	err = userMockRepo.UpdateUsername(ctx, user.ID, "username_new")
	require.NoError(t, err)
	
	user2, err := userMockRepo.GetUserByEmail(ctx, "email3")
	require.NoError(t, err)
	require.Equal(t, user2.Username, "username_new")
	require.Equal(t, user2.Email, "email3")
	require.Equal(t, user2.Password, "password3")
}
