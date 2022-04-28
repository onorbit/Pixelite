package globaldb

import "database/sql"

var gStmtInsertLibrary *sql.Stmt
var gStmtDeleteLibrary *sql.Stmt

type LibraryRow struct {
	RootPath string `db:"root_path"`
}

func initLibraries() error {
	if _, err := gDatabase.Exec("CREATE TABLE IF NOT EXISTS libraries(root_path TEXT PRIMARY KEY)"); err != nil {
		return err
	}

	stmt, err := gDatabase.Prepare("INSERT INTO libraries(root_path) VALUES (?)")
	if err != nil {
		return err
	}
	gStmtInsertLibrary = stmt

	stmt, err = gDatabase.Prepare("DELETE FROM libraries WHERE root_path = ?")
	if err != nil {
		return err
	}
	gStmtDeleteLibrary = stmt

	return nil
}

func InsertLibrary(rootPath string) error {
	_, err := gStmtInsertLibrary.Exec(rootPath)
	return err
}

func DeleteLibrary(rootPath string) error {
	_, err := gStmtDeleteLibrary.Exec(rootPath)
	return err
}

func LoadAllLibraries() ([]LibraryRow, error) {
	ret := []LibraryRow{}
	if err := gDatabase.Select(&ret, "SELECT root_path FROM libraries"); err != nil {
		return nil, err
	}

	return ret, nil
}
