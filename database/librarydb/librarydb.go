package librarydb

import (
	"fmt"
	"math/rand"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type LibraryDB struct {
	libraryID  string
	dbFilePath string
	db         *sqlx.DB
}

func newLibraryDB(dbFilePath string) (*LibraryDB, error) {
	db, err := sqlx.Connect("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	libDB := &LibraryDB{
		dbFilePath: dbFilePath,
		db:         db,
	}

	// initialize tables and stuff if necessary.
	libDB.initMetadata()

	// load or generate ID.
	id, err := libDB.GetMetadata(MetadataKeyLibraryID)
	if err == ErrMetadataNotFound {
		id = fmt.Sprintf("%08x", rand.Uint32())
		libDB.SetMetadata(MetadataKeyLibraryID, id)
	}

	libDB.libraryID = id

	return libDB, nil
}

func (l *LibraryDB) GetLibraryID() string {
	return l.libraryID
}

func (l *LibraryDB) GetDBFilePath() string {
	return l.dbFilePath
}
