package main

import (
	"authentication-service/data"
	"authentication-service/internal/config"
	"authentication-service/internal/driver"
	"authentication-service/internal/repository/dbrepo"
	"authentication-service/internal/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var authPort = os.Getenv("port")

var app config.Config

func main() {
	db := driver.ConnectToDB()
	app = config.Config{
		DB: db,
	}
	data.MigrateUp()

	_ = dbrepo.NewPostgresRepo(app.DB)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", authPort),
		Handler: routes.GetRoutes(),
	}
	fmt.Printf("Starting authentication service on port: %s\n", authPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
