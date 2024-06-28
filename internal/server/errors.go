package server

import (
	"net/http"
)

// Web server errors

func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// API Errors
func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := M{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.Logger.Err(err).Msgf("%s", message)
		w.WriteHeader(500)
	}
}
