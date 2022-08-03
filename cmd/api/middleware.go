package main

import (
	"net/http"
)

func (app *application) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		plainTextToken, err := app.readAuthHeader(r)
		if err != nil {
			app.noAuthorizationHeaderResponse(w, r)
			return
		}

		_, err = app.getValidToken(plainTextToken)
		if err != nil {
			app.failedAuthenticationResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		plaintextToken, err := app.readAuthHeader(r)
		if err != nil {
			app.noAuthorizationHeaderResponse(w, r)
			return
		}

		token, err := app.getValidToken(plaintextToken)
		if err != nil {
			app.failedAuthenticationResponse(w, r)
			return
		}

		user, err := app.models.Tokens.GetUserForToken(token)
		if err != nil {
			app.failedAuthenticationResponse(w, r)
			return
		}

		if user.AccountType != "admin" {
			app.failedAuthenticationResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
