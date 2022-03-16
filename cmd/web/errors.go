package main

import (
	"net/http"
)

// Web server errors

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, http.StatusNotFound)
}
func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	app.clientError(w, http.StatusMethodNotAllowed)
}

// API Errors
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logger.Err(err).Msgf("%s", message)
		w.WriteHeader(500)
	}
}
