package librarydb

type AlbumRow struct {
	ID            int64  `db:"id"`
	Path          string `db:"path"`
	CoverFileName string `db:"cover_filename"`
}

func (l *LibraryDB) initAlbums() error {
	schemaAlbums := `
		CREATE TABLE IF NOT EXISTS albums (
			id INTEGER PRIMARY KEY,
			path TEXT UNIQUE,
			cover_filename TEXT
		)`

	if _, err := l.db.Exec(schemaAlbums); err != nil {
		return err
	}

	return nil
}

func (l *LibraryDB) InsertAlbum(path, coverFileName string) (int64, error) {
	result, err := l.db.Exec("INSERT INTO albums (path, cover_filename) VALUES (?, ?)", path, coverFileName)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (l *LibraryDB) DeleteAlbum(id int64) error {
	_, err := l.db.Exec("DELETE FROM albums WHERE id = ?", id)
	return err
}

func (l *LibraryDB) LoadAllAlbums() ([]AlbumRow, error) {
	ret := []AlbumRow{}
	if err := l.db.Select(&ret, "SELECT id, path, cover_filename FROM albums"); err != nil {
		return nil, err
	}

	return ret, nil
}
