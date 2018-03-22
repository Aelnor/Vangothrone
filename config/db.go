package config

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./vang.db"

func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1) // because of SQLite

	return db, nil
}
