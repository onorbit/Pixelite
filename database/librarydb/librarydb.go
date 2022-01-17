package librarydb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type LibraryDB struct {
	libraryID string
	db        *sql.DB
}
