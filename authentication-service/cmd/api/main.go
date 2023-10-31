package main

import (
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
)

var authPort = os.Getenv("port")

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})
	fmt.Printf("Starting authentication service on port: %s\n", authPort)
	err := http.ListenAndServe(fmt.Sprintf(":%s", authPort), nil)
	if err != nil {
		log.Fatal("cannot start the server")
	}
}
