package repo_test

import (
	"context"
	"log"
	"server/dbTest"
	"server/internal/port"
	"server/internal/repo"
	"testing"
)

var userMockRepo port.UserRepoPort
var chatroomMockRepo port.ChatroomRepoPort
var dbMock *dbTest.DatabaseTest

func TestMain(m *testing.M) {
	db2, err := dbTest.NewDatabaseTest()
	if err != nil {
		log.Fatalf("Something went wrong. Could not connect to the database. %s", err)
	}
	dbMock = db2
	chatroomMockRepo = repo.NewChatroomRepository(dbMock.GetDB())
	userMockRepo = repo.NewUserRepository(dbMock.GetDB())
	m.Run()

	userMockRepo.DeleteUserAll(context.Background())
	chatroomMockRepo.DeleteChatroomAll(context.Background())
	dbMock.Close()
}
