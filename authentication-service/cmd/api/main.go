package main

import (
	"authentication-service/data"
	"authentication-service/internal/config"
	"authentication-service/internal/driver"
	"authentication-service/internal/models"
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
var gRpcPort = os.Getenv("gRpcPort")

var app config.Config

func main() {
	db := driver.ConnectToDB()
	data.MigrateUp()

	dbRepo := dbrepo.NewPostgresRepo(db)

	app = config.Config{
		DB:     dbRepo,
		Models: models.New(),
	}

	go gRPCListen()

	// IDK what to do with server because he is useless
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
