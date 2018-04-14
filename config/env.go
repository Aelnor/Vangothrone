package config

import (
	"database/sql"
)

type Env struct {
	DB *sql.DB
}

func GetStaticPath() string {
	return "/var/vangothrone/"
}
