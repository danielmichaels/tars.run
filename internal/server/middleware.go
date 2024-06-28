package server

import (
	"fmt"
	"net/http"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// If there's a panic, set "Connection: close" on the response.
				// This will tell Go's HTTP server to automatically close the
				// current connection after the response has been sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has a type interface{}, so we
				// use fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. This will log the error using
				// our custom Logger type at the ERROR level and send the client
				// a 500 status.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
