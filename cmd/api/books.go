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
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, headers)
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

func (app *application) getAllReviewsByUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		User string
		data.Filters
	}
	v := validator.NewValidator()
	qs := r.URL.Query()
	input.User = app.readString(qs, "user", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "user", "-id", "-user"}

	data.ValidateFilters(v, input.Filters)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	reviews, metadata, err := app.models.Books.GetAllReviewsByUser(input.User, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"reviews": reviews, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Rating int    `json:"rating"`
		Review string `json:"review"`
		BookID int64  `json:"book_id"`
		UserID int64  `json:"user_id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var review data.Review
	review.Rating = input.Rating
	review.Review = input.Review
	review.BookID = input.BookID
	review.UserID = input.UserID
	v := validator.NewValidator()
	data.ValidateReview(v, &review)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Books.InsertReview(&review)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"review": review}, nil)

}

func (app *application) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}
	review, err := app.models.Books.GetReviewByID(id)
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
		Rating *int    `json:"rating"`
		Review *string `json:"review"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.NewValidator()

	if input.Rating != nil {
		review.Rating = *input.Rating
		v.Check(review.Rating > 0, "rating", "should not be empty")
		v.Check(review.Rating <= 5, "rating", "should not be more than 5")
	}

	if input.Review != nil {
		review.Review = *input.Review
		v.Check(review.Review != "", "review", "should not be empty")
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	review.ID = id
	err = app.models.Books.UpdateReview(review)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"review": review}, nil)
}

func (app *application) getReviewByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}
	review, err := app.models.Books.GetReviewByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"review": review}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamID(r)
	if err != nil {
		app.notfoundResponse(w, r)
		return
	}
	err = app.models.Books.DeleteReview(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.notfoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	message := fmt.Sprintf("review with ID: %d deleted", id)
	err = app.writeJSON(w, http.StatusOK, envelope{"success": message}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
