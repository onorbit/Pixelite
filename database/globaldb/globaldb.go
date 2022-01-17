package globaldb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var gDatabase *sql.DB

func Initialize(dbFilePath string) error {
	database, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return err
	}
	gDatabase = database

	initThumbnails()
	initLibraries()

	return nil
}
