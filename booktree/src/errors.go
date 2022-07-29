package main

import (
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	errMessage := map[string]interface{}{
		"message": message,
	}
	err := app.writeJSON(w, status, errMessage, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
		w.WriteHeader(500)
	}
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.PrintError(err, nil)
	message := "internal server error"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) duplicateUsernameErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "this username is in use"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}

func (app *application) invalidUserErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid user credentials"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) invalidAuthResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
