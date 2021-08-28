package globaldb

import "database/sql"

var gStmtInsertLibrary *sql.Stmt
var gStmtDeleteLibrary *sql.Stmt

type LibraryRow struct {
	ID       string
	RootPath string
	Desc     string
}

func initLibraries() error {
	stmt, err := gDatabase.Prepare("CREATE TABLE IF NOT EXISTS libraries(id TEXT PRIMARY KEY, root_path TEXT, desc TEXT)")
	if err != nil {
		return err
	}
	stmt.Exec()

	stmt, err = gDatabase.Prepare("INSERT INTO libraries(id, root_path, desc) VALUES (?, ?, ?)")
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
	rows, err := gDatabase.Query("SELECT id, root_path, desc FROM libraries")
	// TODO : Query() returns 'no such table' error when the table is newly created. should be handled properly.
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]LibraryRow, 0)
	for rows.Next() {
		var entry LibraryRow

		err := rows.Scan(&entry.ID, &entry.RootPath, &entry.Desc)
		if err != nil {
			return nil, err
		}

		ret = append(ret, entry)
	}

	return ret, nil
}
