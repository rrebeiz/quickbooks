package main

import (
	"fmt"
	"github.com/rrebeiz/quickbooks/internal/data"
	"github.com/rrebeiz/quickbooks/internal/validator"
	"net/http"
	"time"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}

	checkPassword, err := user.Password.CheckPassword(input.Password, app.config.db.pepper)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !checkPassword {
		app.notAuthorizedResponse(w, r)
		return
	}

	token, err := app.models.Tokens.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token.Email = user.Email
	err = app.models.Tokens.InsertToken(token)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user.Token = *token
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) authenticateToken(w http.ResponseWriter, r *http.Request) {
	token, err := app.readAuthHeader(r)
	if err != nil {
		app.noAuthorizationHeaderResponse(w, r)
		return
	}

	tkn, err := app.models.Tokens.GetByToken(*token)
	if err != nil {
		app.notAuthorizedResponse(w, r)
		return
	}

	if tkn.Expiry.Before(time.Now()) {
		app.notAuthorizedResponse(w, r)
		return
	}

	user, err := app.models.Tokens.GetUserForToken(tkn)
	if err != nil {
		app.notAuthorizedResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	token, err := app.readAuthHeader(r)
	if err != nil {
		app.noAuthorizationHeaderResponse(w, r)
		return
	}
	tkn, err := app.models.Tokens.GetByToken(*token)
	if err != nil {
		app.notAuthorizedResponse(w, r)
		return
	}
	if tkn.Expiry.Before(time.Now()) {
		app.notAuthorizedResponse(w, r)
	}

	err = app.models.Tokens.DeleteToken(tkn.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	message := "token destroyed"
	err = app.writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) adminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}
	err = app.models.Tokens.DeleteToken(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	message := fmt.Sprintf("token with id %d destroyed", id)
	err = app.writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
