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

	if err = initThumbnails(); err != nil {
		return err
	}

	if err = initLibraries(); err != nil {
		return err
	}

	return nil
}
