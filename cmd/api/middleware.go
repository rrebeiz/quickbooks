package main

import (
	"errors"
	"github.com/rrebeiz/quickbooks/internal/data"
	"net/http"
)

func (app *application) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		plainTextToken, err := app.readAuthHeader(r)
		if err != nil {
			switch {
			case errors.Is(err, ErrNoAuthHeader):
				app.noAuthorizationHeaderResponse(w, r)
			case errors.Is(err, data.ErrNoRecordFound):
				app.notAuthorizedResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		_, err = app.getValidToken(plainTextToken)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRecordFound):
				app.notAuthorizedResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		plaintextToken, err := app.readAuthHeader(r)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRecordFound):
				app.notAuthorizedResponse(w, r)
			case errors.Is(err, ErrNoAuthHeader):
				app.noAuthorizationHeaderResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		token, err := app.getValidToken(plaintextToken)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRecordFound):
				app.notAuthorizedResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		user, err := app.models.Tokens.GetUserForToken(token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrNoRecordFound):
				app.notAuthorizedResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		if user.AccountType != "admin" {
			app.notAuthorizedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
