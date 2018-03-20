package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var dbPath = "./vang.db"
var db *sql.DB

func initDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(1) // because of SQLite

	return nil
}
