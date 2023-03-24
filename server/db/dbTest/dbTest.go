package dbTest

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DatabaseTest struct {
	db *sql.DB
}

func NewDatabaseTest() (*DatabaseTest, error) {
	db, err := sql.Open("postgres", "postgresql://root:password@localhost:5433/go-chat-test?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &DatabaseTest{db: db}, nil
}

func (d *DatabaseTest) Close() error {
	return d.db.Close()
}

func (d *DatabaseTest) GetDB() *sql.DB {
	return d.db
}
