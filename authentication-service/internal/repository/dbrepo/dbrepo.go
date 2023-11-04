package dbrepo

import (
	"authentication-service/internal/repository"
	"database/sql"
)

type postgresDBRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB: conn,
	}
}
