package config

import (
	"database/sql"
)

type Config struct {
	DB *sql.DB
}
