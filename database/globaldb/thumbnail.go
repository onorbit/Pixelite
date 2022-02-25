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
			last_access_timestamp INTEGER
		)`

	if _, err := gDatabase.Exec(schemaThumbnailedAlbums); err != nil {
		return err
	}

	return nil
}

func RegisterThumbnail(imagePath, thumbnailPath string, thumbnailedAlbumID int64) error {
	_, err := gDatabase.Exec("INSERT INTO thumbnails (image_path, thumbnail_path, thumbnailed_album_id) VALUES (?, ?, ?)", imagePath, thumbnailPath, thumbnailedAlbumID)
	return err
}

func InsertThumbnailedAlbum(libraryID, albumID string, lastAccessTimestamp time.Time) (int64, error) {
	result, err := gDatabase.Exec("INSERT INTO thumbnailed_albums (library_id, album_id, last_access_timestamp) VALUES (?, ?, ?)", libraryID, albumID, lastAccessTimestamp.Unix())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
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
	if err := gDatabase.Select(&ret, "SELECT id, library_id, album_id, last_access_timestamp FROM thumbnailed_albums"); err != nil {
		return nil, err
	}

	return ret, nil
}
