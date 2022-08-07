package main

import (
	"errors"
	"fmt"
	"github.com/rrebeiz/quickbooks/internal/data"
	"github.com/rrebeiz/quickbooks/internal/validator"
	"net/http"
)

func (app *application) getAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		data.Filters
	}
	v := validator.NewValidator()

	qs := r.URL.Query()
	input.Title = app.readString(qs, "title", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "publication_year", "-id", "-title", "-publication_year"}

	data.ValidateFilters(v, input.Filters)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	books, metadata, err := app.models.Books.GetAll(input.Title, input.Filters)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"books": books, "metadata": metadata}, nil)
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
		Slug string
	}
	qs := r.URL.Query()
	input.Slug = app.readString(qs, "slug", "")

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

func (app *application) getAllAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Author string
		data.Filters
	}
	v := validator.NewValidator()

	qs := r.URL.Query()
	input.Author = app.readString(qs, "author_name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "author_name", "-id", "-author_name"}

	data.ValidateFilters(v, input.Filters)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	authors, metadata, err := app.models.Books.GetAllAuthors(input.Author, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"authors": authors, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
