package main

import (
	"fmt"
	mh "github.com/gorilla/handlers"
	"log"
	"net/http"
	"os"
	"shorty"
	"shorty/adapters"
	"shorty/handlers"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	err := startServer()
	if err != nil {
		log.Fatalln("failed to start server", err)
	}

	return nil
}
func startServer() error {
	log.Println("Starting server")
	conf := shorty.AppConfig()

	database := conf.Db.DbName
	adapters.InitialMigrations(database)

	s := handlers.NewServer()
	http.Handle("/", mh.LoggingHandler(os.Stdout, s.Router()))

	port := conf.Server.Port
	log.Printf("Listening on port: %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

	return nil
}
