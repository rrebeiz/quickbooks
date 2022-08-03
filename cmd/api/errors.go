package main

import (
	"fmt"
	"net/http"
)

func (app *application) printError(err error) {
	app.errorLog.Println(err)
}
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {

	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.errorLog.Println(err, r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.printError(err)
	message := "the server encountered a problem."
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) notfoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the requested method %s is not allowed on this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) failedAuthenticationResponse(w http.ResponseWriter, r *http.Request) {
	message := "authentication failed"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) noAuthorizationHeaderResponse(w http.ResponseWriter, r *http.Request) {
	message := "no authorization header received"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}

func (app *application) duplicateEmailResponse(w http.ResponseWriter, r *http.Request) {
	message := "email address already taken"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}
