package globaldb

import "database/sql"

var gStmtInsertThumbnail *sql.Stmt

type ThumbnailRow struct {
	ImagePath     string `db:"image_path"`
	ThumbnailPath string `db:"thumbnail_path"`
}

func initThumbnails() error {
	if _, err := gDatabase.Exec("CREATE TABLE IF NOT EXISTS thumbnails(image_path TEXT PRIMARY KEY, thumbnail_path TEXT)"); err != nil {
		return err
	}

	stmt, err := gDatabase.Prepare("INSERT INTO thumbnails (image_path, thumbnail_path) VALUES (?, ?)")
	if err != nil {
		return err
	}
	gStmtInsertThumbnail = stmt

	return nil
}

func RegisterThumbnail(imagePath, thumbnailPath string) error {
	_, err := gStmtInsertThumbnail.Exec(imagePath, thumbnailPath)
	return err
}

func LoadAllThumbnails() ([]ThumbnailRow, error) {
	ret := []ThumbnailRow{}
	if err := gDatabase.Select(&ret, "SELECT image_path, thumbnail_path FROM thumbnails"); err != nil {
		return nil, err
	}

	return ret, nil
}
