package globaldb

import "database/sql"

var gStmtInsertThumbnail *sql.Stmt

type ThumbnailRow struct {
	ImagePath     string
	ThumbnailPath string
}

func initThumbnails() error {
	stmt, err := gDatabase.Prepare("CREATE TABLE IF NOT EXISTS thumbnails(image_path TEXT PRIMARY KEY, thumbnail_path TEXT)")
	if err != nil {
		return err
	}
	stmt.Exec()

	stmt, err = gDatabase.Prepare("INSERT INTO thumbnails (image_path, thumbnail_path) VALUES (?, ?)")
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
	rows, err := gDatabase.Query("SELECT image_path, thumbnail_path FROM thumbnails")
	// TODO : Query() returns 'no such table' error when the table is newly created. should be handled properly.
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]ThumbnailRow, 0)
	for rows.Next() {
		var entry ThumbnailRow

		err = rows.Scan(&entry.ImagePath, &entry.ThumbnailPath)
		if err != nil {
			return nil, err
		}

		ret = append(ret, entry)
	}

	return ret, nil
}
