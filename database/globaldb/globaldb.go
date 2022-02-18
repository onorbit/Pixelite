package globaldb

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var gDatabase *sqlx.DB

func Initialize(dbFilePath string) error {
	database, err := sqlx.Connect("sqlite3", dbFilePath)
	if err != nil {
		return err
	}

	gDatabase = database

	initThumbnails()
	initLibraries()

	return nil
}
