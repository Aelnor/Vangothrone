package models

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

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func InitUsersTable(db *sql.DB) error {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Users(login, name, password, is_admin)"); err != nil {
		return err
	}

	return AddUser(db, "eiden", "AC", "pass", true)
}

func checkCredentials(db *sql.DB, login string, password string) (*User, error) {
	row := db.QueryRow("SELECT rowid, login, name, is_admin FROM Users WHERE login=? AND password=?", login, GetMD5Hash(password))

	u := new(User)
	err := row.Scan(&u.Id, &u.Login, &u.Name, &u.IsAdmin)

	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("Login or password are incorrent")
	case err != nil:
		return nil, err
	}

	return u, nil
}

func AddUser(db *sql.DB, login string, name string, password string, isAdmin bool) error {
	statement, err := db.Prepare("INSERT INTO Users(login, name, password, is_admin) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(login, name, GetMD5Hash(password), isAdmin)

	return err
}
