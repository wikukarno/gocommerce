package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3" // Menggunakan underscore untuk import yang hanya efek sampingnya
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./gocommerce.db")
	if err != nil {
		panic(err)
	}

	createUsersTable()
}

func createUsersTable() {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		password TEXT
	);
	`
	_, err := DB.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}
}