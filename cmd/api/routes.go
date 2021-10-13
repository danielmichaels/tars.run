package main

import (
	mh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	//router.HandlerFunc(http.MethodGet, "/v1/links", app.healthcheckHandler)
	router.HandleFunc("/v1/healthcheck", app.healthcheckHandler).Methods(http.MethodGet)

	return mh.LoggingHandler(os.Stdout, app.recoverPanic(app.enableCORS(router)))
}
