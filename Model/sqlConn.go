package Model

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func checkErrorDB(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error connecting to database:", err)
		panic(err)
	}
}

type DB struct {
	db *sql.DB
}

func (d *DB) connect(path string) {
	db, err := sql.Open("sqlite3", path)
	checkErrorDB(err)
	checkErrorDB(db.Ping())
	d.db = db
}

func (d DB) stop() {
	d.db.Close()
}

func (d DB) query(querystring string, model TableType) Table {
	rows, err := d.db.Query(querystring)
	checkErrorDB(err)
	return model.query(rows)
}

func (d DB) count(tableName string) (cnt int) {
	queryString := `select count(*) from '` + tableName + `'`
	err := database.db.QueryRow(queryString, 1).Scan(&cnt)
	checkErrorDB(err)
	return
}
