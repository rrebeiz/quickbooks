package main

import (
	"errors"
	"fmt"
	"github.com/rrebeiz/quickbooks/internal/data"
	"github.com/rrebeiz/quickbooks/internal/validator"
	"net/http"
)

func (app *application) getAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := app.models.Books.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}
	book, err := app.models.Books.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getBookBySlugHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Slug string `json:"slug"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	book, err := app.models.Books.GetBySlug(input.Slug)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title           string `json:"title"`
		AuthorID        int    `json:"author_id"`
		PublicationYear int    `json:"publication_year"`
		Description     string `json:"description"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var book data.Book
	book.Title = input.Title
	book.AuthorID = input.AuthorID
	book.PublicationYear = input.PublicationYear
	book.Description = input.Description
	book.Slug = book.Title

	v := validator.NewValidator()
	data.ValidateBook(v, &book)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Insert(&book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}
	book, err := app.models.Books.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Title           *string  `json:"title"`
		AuthorID        *int     `json:"author_id"`
		PublicationYear *int     `json:"publication_year"`
		Description     *string  `json:"description"`
		Genres          []string `json:"genres"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.NewValidator()

	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.AuthorID != nil {
		book.AuthorID = *input.AuthorID
	}
	if input.PublicationYear != nil {
		book.PublicationYear = *input.PublicationYear
	}
	if input.Description != nil {
		book.Description = *input.Description
	}
	if input.Genres != nil {
		book.Genres = input.Genres
		v.Check(validator.Unique(book.Genres), "genres", "values must be unique")
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		var genreList = []string{"Science Fiction", "Fantasy", "Romance", "Thriller", "Mystery", "Horror", "Classic", "Self-help"}
		msg := fmt.Sprintf("please use the following genres %s", genreList)
		for i := range book.Genres {
			v.Check(validator.PermittedValue(book.Genres[i], genreList...), "genres", msg)
			if !v.Valid() {
				app.failedValidationResponse(w, r, v.Errors)
				return
			}
		}
	}
	data.ValidateBook(v, book)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Update(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}
	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	msg := fmt.Sprintf("book with ID %d deleted.", id)
	err = app.writeJSON(w, http.StatusOK, envelope{"message": msg}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
