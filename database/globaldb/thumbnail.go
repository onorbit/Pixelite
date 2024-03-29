package globaldb

import (
	"time"
)

type ThumbnailRow struct {
	ImagePath          string `db:"image_path"`
	ThumbnailPath      string `db:"thumbnail_path"`
	ThumbnailedAlbumID int64  `db:"thumbnailed_album_id"`
}

type ThumbnailedAlbumRow struct {
	ID                  int64  `db:"id"`
	LibraryID           string `db:"library_id"`
	AlbumID             string `db:"album_id"`
	CreateTimestamp     int64  `db:"create_timestamp"`
	LastAccessTimestamp int64  `db:"last_access_timestamp"`
}

type AlbumCoverRow struct {
	LibraryID           string `db:"library_id"`
	AlbumID             string `db:"album_id"`
	CoverPath           string `db:"cover_path"`
	LastAccessTimestamp int64  `db:"last_access_timestamp"`
}

func initThumbnails() error {
	// thumbnails table and index.
	schemaThumbnails := `
		CREATE TABLE IF NOT EXISTS thumbnails (
			image_path TEXT PRIMARY KEY,
			thumbnail_path TEXT,
			thumbnailed_album_id INTEGER
		)`

	if _, err := gDatabase.Exec(schemaThumbnails); err != nil {
		return err
	}

	if _, err := gDatabase.Exec("CREATE INDEX IF NOT EXISTS thumbnails__thumbnailed_album_id ON thumbnails (thumbnailed_album_id)"); err != nil {
		return err
	}

	// thumbnailed_albums table.
	schemaThumbnailedAlbums := `
		CREATE TABLE IF NOT EXISTS thumbnailed_albums (
			id INTEGER PRIMARY KEY,
			library_id TEXT,
			album_id TEXT,
			create_timestamp INTEGER,
			last_access_timestamp INTEGER
		)`

	if _, err := gDatabase.Exec(schemaThumbnailedAlbums); err != nil {
		return err
	}

	// album_covers table.
	schemaAlbumCovers := `
		CREATE TABLE IF NOT EXISTS album_covers (
			library_id TEXT,
			album_id TEXT,
			cover_path TEXT,
			last_access_timestamp INTEGER,
			PRIMARY KEY (library_id, album_id)
		)`

	if _, err := gDatabase.Exec(schemaAlbumCovers); err != nil {
		return err
	}

	return nil
}

func InsertThumbnail(imagePath, thumbnailPath string, thumbnailedAlbumID int64) error {
	_, err := gDatabase.Exec("INSERT INTO thumbnails (image_path, thumbnail_path, thumbnailed_album_id) VALUES (?, ?, ?)", imagePath, thumbnailPath, thumbnailedAlbumID)
	return err
}

func InsertThumbnailedAlbum(libraryID, albumID string) (int64, error) {
	currTime := time.Now()
	result, err := gDatabase.Exec("INSERT INTO thumbnailed_albums (library_id, album_id, create_timestamp, last_access_timestamp) VALUES (?, ?, ?, ?)", libraryID, albumID, currTime.Unix(), currTime.Unix())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func InsertAlbumCover(libraryID, albumID, coverPath string) error {
	currTime := time.Now()
	_, err := gDatabase.Exec("INSERT INTO album_covers (library_id, album_id, cover_path, last_access_timestamp) values (?, ?, ?, ?)", libraryID, albumID, coverPath, currTime.Unix())
	return err
}

// DeleteThumbnailedAlbum deletes both thumbnailed_albums table entry and belonging thumbnails table entries.
func DeleteThumbnailedAlbum(thumbnailedAlbumID int64) (err error) {
	tx, err := gDatabase.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = r.(error)
		}
	}()

	gDatabase.MustExec("DELETE FROM thumbnails WHERE thumbnailed_album_id = ?", thumbnailedAlbumID)
	gDatabase.MustExec("DELETE FROM thumbnailed_albums WHERE id = ?", thumbnailedAlbumID)

	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return
}

func UpdateThumbnailedAlbumAccessTimestamp(thumbnailedAlbumID int64, timestamp time.Time) error {
	_, err := gDatabase.Exec("UPDATE thumbnailed_albums SET last_access_timestamp = ? WHERE id = ?", timestamp.Unix(), thumbnailedAlbumID)
	return err
}

func LoadAllThumbnails() ([]ThumbnailRow, error) {
	ret := []ThumbnailRow{}
	if err := gDatabase.Select(&ret, "SELECT image_path, thumbnail_path, thumbnailed_album_id FROM thumbnails"); err != nil {
		return nil, err
	}

	return ret, nil
}

func LoadAllThumbnailedAlbums() ([]ThumbnailedAlbumRow, error) {
	ret := []ThumbnailedAlbumRow{}
	if err := gDatabase.Select(&ret, "SELECT id, library_id, album_id, create_timestamp, last_access_timestamp FROM thumbnailed_albums"); err != nil {
		return nil, err
	}

	return ret, nil
}

func LoadAllAlbumCovers() ([]AlbumCoverRow, error) {
	ret := []AlbumCoverRow{}
	if err := gDatabase.Select(&ret, "SELECT library_id, album_id, cover_path, last_access_timestamp FROM album_covers"); err != nil {
		return nil, err
	}

	return ret, nil
}
