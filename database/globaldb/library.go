package globaldb

import "database/sql"

var gStmtInsertLibrary *sql.Stmt
var gStmtDeleteLibrary *sql.Stmt

type LibraryRow struct {
	ID       string `db:"id"`
	RootPath string `db:"root_path"`
	Desc     string `db:"desc"`
}

func initLibraries() error {
	if _, err := gDatabase.Exec("CREATE TABLE IF NOT EXISTS libraries(id TEXT PRIMARY KEY, root_path TEXT, desc TEXT)"); err != nil {
		return err
	}

	stmt, err := gDatabase.Prepare("INSERT INTO libraries(id, root_path, desc) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	gStmtInsertLibrary = stmt

	stmt, err = gDatabase.Prepare("DELETE FROM libraries WHERE id = ?")
	if err != nil {
		return err
	}
	gStmtDeleteLibrary = stmt

	return nil
}

func InsertLibrary(id, rootPath, desc string) error {
	_, err := gStmtInsertLibrary.Exec(id, rootPath, desc)
	return err
}

func DeleteLibrary(id string) error {
	_, err := gStmtDeleteLibrary.Exec(id)
	return err
}

func LoadAllLibraries() ([]LibraryRow, error) {
	ret := []LibraryRow{}
	if err := gDatabase.Select(&ret, "SELECT id, root_path, desc FROM libraries"); err != nil {
		return nil, err
	}

	return ret, nil
}
