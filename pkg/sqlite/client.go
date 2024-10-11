package sqlite

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"time"
)

type Client interface {
	OpenConnection() (*sql.DB, error)
	CloseConnection(db *sql.DB) error
}

type client struct {
}

func (self *client) OpenConnection() (db *sql.DB, err error) {

	if db, err = sql.Open("sqlite", "./data/sqlite-database.db"); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(1 * time.Minute)

	return db, nil
}

func (self *client) CloseConnection(db *sql.DB) error {
	return db.Close()
}

func NewClient(initDumpData bool) Client {

	if initDumpData {
		CreateTable()
		Insert()
	}

	return &client{}
}
