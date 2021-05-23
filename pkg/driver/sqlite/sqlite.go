package sqlite

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

type ResultCalc struct {
	Id      int
	A       int
	B       int
	C       int
	Delta   int
	X1      float64
	X2      float64
	IsValid bool
}

func NewSQLite() (SQLite, error) {

	if _, err := os.Stat("./database.db"); os.IsNotExist(err) {
		file, err := os.Create("database.db")
		if err != nil {
			return SQLite{}, err
		}
		file.Close()
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return SQLite{}, err
	}

	var sqlite SQLite
	sqlite.db = db
	err = sqlite.createTable()
	if err != nil {
		return SQLite{}, err
	}
	return sqlite, nil

}

func (db *SQLite) createTable() error {

	sqlQuery := `CREATE TABLE IF NOT EXISTS calculations (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"a" integer,
		"b" integer,
		"c" integer,
		"delta" integer,
		"x1" text,
		"x2" text,
		"is_valid" boolean
		)`

	query, err := db.db.Prepare(sqlQuery)
	if err != nil {
		return err
	}

	_, err = query.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (db *SQLite) ReadCalculation(a, b, c int) (*ResultCalc, error) {
	sqlQuery := `SELECT * FROM calculations WHERE a=? AND b=? AND c=?`
	query, err := db.db.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}

	result := query.QueryRow(a, b, c)

	var resCalc ResultCalc
	err = result.Scan(&resCalc.Id, &resCalc.A, &resCalc.B, &resCalc.C, &resCalc.Delta, &resCalc.X1, &resCalc.X2, &resCalc.IsValid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &resCalc, err
}

func (db *SQLite) PutCalculation(calculation ResultCalc) error {
	sqlQuery := `INSERT INTO calculations (a, b, c, delta, x1, x2, is_valid) VALUES (?, ?, ?, ?, ?, ?, ?)`
	query, err := db.db.Prepare(sqlQuery)
	if err != nil {
		return err
	}

	_, err = query.Exec(
		calculation.A,
		calculation.B,
		calculation.C,
		calculation.Delta,
		calculation.X1,
		calculation.X2,
		calculation.IsValid,
	)
	if err != nil {
		return err
	}
	return nil
}
