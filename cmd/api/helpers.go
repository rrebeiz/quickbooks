package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/rrebeiz/quickbooks/internal/data"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type envelope map[string]any

var (
	ErrNoAuthHeader = errors.New("no authorization header provided")
)

func (app *application) readParamID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParamFromCtx(r.Context(), "id"), 10, 64)
	if err != nil || id < 1 {
		return 0, err
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	var output []byte
	if app.config.env == "development" {
		js, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return err
		}
		output = js
	} else {
		js, err := json.Marshal(data)
		if err != nil {
			return err
		}
		output = js
	}

	for key, value := range headers {
		w.Header()[key] = value
	}
	output = append(output, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(output)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}
	return nil
}

func (app *application) readAuthHeader(r *http.Request) (*string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, ErrNoAuthHeader
	}
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, ErrNoAuthHeader
	}
	token := headerParts[1]
	if len(token) != 26 {
		return nil, ErrNoAuthHeader
	}

	tkn, err := app.models.Tokens.GetByToken(token)
	if err != nil {
		return nil, err
	}
	return &tkn.Token, nil
}

func (app *application) getValidToken(plainTextToken *string) (*data.Token, error) {

	token, err := app.models.Tokens.GetByToken(*plainTextToken)
	if err != nil {
		return nil, err
	}

	if token.Expiry.Before(time.Now()) {
		return nil, err
	}
	return token, nil
}
