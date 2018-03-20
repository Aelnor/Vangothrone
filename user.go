package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
)

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Login   string `json:"login"`
	IsAdmin bool   `json:"isAdmin"`
}

var CurrentUser User

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func initUsersTable(db *sql.DB) error {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Users(login, name, password, is_admin)"); err != nil {
		return err
	}

	return addUser(db, "eiden", "AC", "pass", true)
}

func checkCredentials(db *sql.DB, login string, password string) error {
	row := db.QueryRow("SELECT rowid, login, name, is_admin FROM Users WHERE login=? AND password=?", login, GetMD5Hash(password))

	err := row.Scan(&CurrentUser.Id, &CurrentUser.Login, &CurrentUser.Name, &CurrentUser.IsAdmin)

	switch {
	case err == sql.ErrNoRows:
		return fmt.Errorf("Login or password are incorrent")
	case err != nil:
		return err
	}

	return nil
}

func addUser(db *sql.DB, login string, name string, password string, isAdmin bool) error {
	statement, err := db.Prepare("INSERT INTO Users(login, name, password, is_admin) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(login, name, GetMD5Hash(password), isAdmin)

	return err
}

func fillUser(username string, password string) error {
	var err error
	err = checkCredentials(db, username, password)
	return err
}
