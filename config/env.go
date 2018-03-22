package config

import (
	"database/sql"

	"github.com/aelnor/vangothrone/models"
)

type Env struct {
	DB          *sql.DB
	CurrentUser models.User
}
