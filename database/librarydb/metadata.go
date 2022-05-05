package librarydb

import (
	"database/sql"
	"errors"
)

type MetadataRow struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

const MetadataKeyLibraryID = "id"
const MetadataKeyLibraryTitle = "title"

var ErrMetadataNotFound = errors.New("metadata with designated key not found")

func (l *LibraryDB) initMetadata() error {
	if _, err := l.db.Exec("CREATE TABLE IF NOT EXISTS metadata(key TEXT PRIMARY KEY, value TEXT)"); err != nil {
		return err
	}

	return nil
}

func (l *LibraryDB) GetMetadata(key string) (string, error) {
	ret := MetadataRow{}
	if err := l.db.Get(&ret, "SELECT key, value FROM metadata WHERE key = $1", key); err != nil {
		if err == sql.ErrNoRows {
			err = ErrMetadataNotFound
		}
		return "", err
	}

	return ret.Value, nil
}

func (l *LibraryDB) SetMetadata(key, value string) (err error) {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = r.(error)
		}
	}()

	// try update first.
	result, err := l.db.Exec("UPDATE metadata SET value = ? WHERE key = ?", value, key)
	if err != nil {
		panic(err)
	}
	nAffected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if nAffected == 1 {
		tx.Commit()
		return nil
	}

	// try insert.
	_, err = l.db.Exec("INSERT INTO metadata(key, value) VALUES (?, ?)", key, value)
	if err != nil {
		panic(err)
	}

	tx.Commit()
	return nil
}
