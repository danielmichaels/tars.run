package main

import (
	mh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	router.HandleFunc("/v1/healthcheck", app.healthcheckHandler).Methods(http.MethodGet)
	router.HandleFunc("/v1/links/{hash}", app.showLinkHandler()).Methods(http.MethodGet)
	router.HandleFunc("/v1/links", app.createLinkHandler()).Methods(http.MethodPost)
	router.HandleFunc("/v1/links/{hash}/analytics", app.showLinkAnalyticsHandler()).Methods(http.MethodGet)

	return mh.LoggingHandler(os.Stdout, app.recoverPanic(app.enableCORS(app.rateLimit(router))))
}
