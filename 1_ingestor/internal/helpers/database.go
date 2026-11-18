package helpers

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver (cgo: needs gcc to build)
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

func OpenDB(databasePath string) error {
	var err error
	// _busy_timeout: if another program (e.g. a SQLite GUI) holds the file, wait up to 5s instead of failing immediately.
	DB, err = sql.Open("sqlite3", "file:"+databasePath+"?_busy_timeout=5000")
	if err != nil {
		return err
	}

	// One connection at a time as SQLite accepts a single writer at a time
	DB.SetMaxOpenConns(1)

	return DB.Ping()
}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		logrus.Errorf("error closing db : %s", err.Error())
	}
}
