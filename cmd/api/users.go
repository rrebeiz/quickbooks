package main

import (
	"errors"
	"fmt"
	"github.com/rrebeiz/quickbooks/internal/data"
	"github.com/rrebeiz/quickbooks/internal/validator"
	"net/http"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var user data.User
	user.Name = input.Name
	user.Email = input.Email
	err = user.Password.HashPassword(input.Password, app.config.db.pepper)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	data.ValidateUser(v, &user)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.duplicateEmailResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil || id < 1 {
		app.notfoundResponse(w, r)
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	v := validator.NewValidator()
	if input.Name != nil {
		data.ValidateName(v, *input.Name)
		user.Name = *input.Name
	}
	if input.Email != nil {
		data.ValidateEmail(v, *input.Email)
		user.Email = *input.Email
	}
	if input.Password != nil {
		data.ValidatePassword(v, *input.Password)
		user.Password.Plaintext = input.Password
		err := user.Password.HashPassword(*input.Password, app.config.db.pepper)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			app.duplicateEmailResponse(w, r)
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil || id < 1 {
		app.notfoundResponse(w, r)
		return
	}
	err = app.models.Users.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	message := fmt.Sprintf("user with the ID: %d has been deleted!", id)
	err = app.writeJSON(w, http.StatusOK, envelope{"message": message}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil || id < 1 {
		app.notfoundResponse(w, r)
		return
	}
	user, err := app.models.Users.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllAuthenticatedUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.GetAllLoggedIn()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
