package models

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
)

type User struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Login   string `json:"login"`
	IsAdmin bool   `json:"isAdmin"`
}

const (
	SELECT_ALL = "SELECT rowid, login, name, is_admin FROM Users"
)

var cachedUsers []*User
var usersMx sync.Mutex

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func InitUsersTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS Users(login, name, password, is_admin)")
	return err
}

func LoadUser(db *sql.DB, login string, password string) (*User, error) {
	row := db.QueryRow("SELECT rowid, login, name, is_admin FROM Users WHERE login=? AND password=?", strings.ToLower(login), password)

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

func CheckCredentials(db *sql.DB, login string, password string) (*User, error) {
	return LoadUser(db, login, GetMD5Hash(password))
}

func AddUser(db *sql.DB, login string, name string, password string, isAdmin bool) error {
	statement, err := db.Prepare("INSERT INTO Users(login, name, password, is_admin) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(strings.ToLower(login), name, GetMD5Hash(password), isAdmin)
	if err == nil {
		invalidateUsersCache()
	}

	return err
}

func LoadUsers(db *sql.DB) ([]*User, error) {
	users := cachedUsers
	if users != nil {
		return users, nil
	}
	usersMx.Lock()
	defer usersMx.Unlock()
	if cachedUsers != nil {
		return cachedUsers, nil
	}
	rows, err := db.Query(SELECT_ALL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users = make([]*User, 0)

	for rows.Next() {
		var id int64
		var login, name string
		var isAdmin bool
		if err := rows.Scan(&id, &login, &name, &isAdmin); err != nil {
			return nil, err
		}

		users = append(users, &User{Id: id, Login: login, Name: name, IsAdmin: isAdmin})
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	cachedUsers = users

	return users, nil
}

func invalidateUsersCache() {
	usersMx.Lock()
	defer usersMx.Unlock()
	cachedUsers = nil
}
