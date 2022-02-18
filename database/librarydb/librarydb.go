package librarydb

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type LibraryDB struct {
	libraryID string
	db        *sqlx.DB
}
