package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mozillazg/go-slugify"
	"github.com/rrebeiz/quickbooks/internal/validator"
	"strconv"
	"strings"
	"time"
)

type Book struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	AuthorID        int       `json:"author_id"`
	PublicationYear int       `json:"publication_year"`
	Slug            string    `json:"slug"`
	Author          Author    `json:"author"`
	Description     string    `json:"description"`
	Genres          []string  `json:"genres"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}
type Author struct {
	ID         int64     `json:"id"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
	Version    int       `json:"version"`
}

type Genre struct {
	ID        int64     `json:"id"`
	GenreName string    `json:"genre_name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Books interface {
	GetAll(genreIDs ...int) ([]*Book, error)
	GetByID(id int64) (*Book, error)
	GetBySlug(slug string) (*Book, error)
	Insert(book *Book) error
	Update(book *Book) error
	Delete(id int64) error
}

type BookModel struct {
	DB *sql.DB
}

func NewBookModel(db *sql.DB) BookModel {
	return BookModel{DB: db}
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Title != "", "title", "should not be empty")
	v.Check(book.Description != "", "description", "should not be empty")
	v.Check(book.PublicationYear > 0, "publication_year", "should not be empty")
}

func (b BookModel) GetAll(genreIDs ...int) ([]*Book, error) {
	where := ""
	if len(genreIDs) > 0 {
		var IDs []string
		for _, x := range genreIDs {
			IDs = append(IDs, strconv.Itoa(x))
		}
		where = fmt.Sprintf("where b.id in (%s)", strings.Join(IDs, ","))
	}
	query := fmt.Sprintf(`select b.id, b.title, b.author_id, b.publication_year, b.slug, b.description, b.created_at, 
						b.updated_at, a.id, a.author_name, a.created_at, a.updated_at, a.version from books b left join authors a on (b.author_id = a.id) %s order by b.title`, where)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []*Book

	rows, err := b.DB.QueryContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.AuthorID, &book.PublicationYear, &book.Slug, &book.Description, &book.CreatedAt,
			&book.UpdatedAt, &book.Author.ID, &book.Author.AuthorName, &book.Author.CreatedAt, &book.Author.UpdatedAt, &book.Author.Version)
		if err != nil {
			return nil, err
		}

		genres, err := b.genresByBook(book.ID)
		if err != nil {
			return nil, err
		}
		for i := range genres {
			book.Genres = append(book.Genres, genres[i].GenreName)
		}
		books = append(books, &book)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (b BookModel) GetByID(id int64) (*Book, error) {
	query := `select b.id, b.title, b.author_id, b.publication_year, b.slug, b.description, b.created_at, b.updated_at, a.id, 
       a.author_name, a.created_at, a.updated_at, a.version from books b left join authors a on (b.author_id = a.id) where b.id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var book Book
	err := b.DB.QueryRowContext(ctx, query, id).Scan(&book.ID, &book.Title, &book.AuthorID, &book.PublicationYear, &book.Slug,
		&book.Description, &book.CreatedAt, &book.UpdatedAt, &book.Author.ID, &book.Author.AuthorName, &book.Author.CreatedAt, &book.Author.UpdatedAt, &book.Author.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	genres, err := b.genresByBook(book.ID)
	if err != nil {
		return nil, err
	}
	var bookGenres []string
	for _, x := range genres {
		bookGenres = append(bookGenres, x.GenreName)
	}
	book.Genres = bookGenres
	return &book, nil
}

func (b BookModel) GetBySlug(slug string) (*Book, error) {
	query := `select b.id, b.title, b.author_id, b.publication_year, b.slug, b.description, b.created_at, b.updated_at, a.id, 
              a.author_name, a.created_at, a.updated_at, a.version from books b left join authors a on (b.author_id = a.id) where b.slug = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var book Book
	err := b.DB.QueryRowContext(ctx, query, slug).Scan(&book.ID, &book.Title, &book.AuthorID, &book.PublicationYear,
		&book.Slug, &book.Description, &book.CreatedAt, &book.UpdatedAt, &book.Author.ID, &book.Author.AuthorName, &book.Author.CreatedAt, &book.Author.UpdatedAt, &book.Author.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	genres, err := b.genresByBook(book.ID)
	if err != nil {
		return nil, err
	}
	var bookGenres []string
	for _, x := range genres {
		bookGenres = append(bookGenres, x.GenreName)
	}
	book.Genres = bookGenres
	return &book, nil
}

func (b BookModel) Insert(book *Book) error {
	query := `insert into books (title, author_id, publication_year, slug, description) values ($1, $2, $3, $4, $5) returning id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{book.Title, book.AuthorID, book.PublicationYear, slugify.Slugify(book.Title), book.Description}
	return b.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID)
}

func (b BookModel) Update(book *Book) error {
	query := `update books set title = $1, author_id = $2, publication_year = $3, slug = $4, description = $5, updated_at = now() where id = $6 returning updated_at`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{book.Title, book.AuthorID, book.PublicationYear, slugify.Slugify(book.Title), book.Description, book.ID}
	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&book.UpdatedAt)
	if err != nil {
		return err
	}

	if len(book.Genres) > 0 {
		query = `delete from books_genres where book_id = $1`
		_, err := b.DB.ExecContext(ctx, query, book.ID)
		if err != nil {
			return fmt.Errorf("failed to delete genres %s", err.Error())
		}
		for i := range book.Genres {
			query = `select id from genres where genre_name = $1`
			var genre Genre
			err = b.DB.QueryRowContext(ctx, query, book.Genres[i]).Scan(&genre.ID)
			if err != nil {
				return err
			}
			query = `insert into books_genres (book_id, genre_id) values ($1, $2)`
			_, err := b.DB.ExecContext(ctx, query, book.ID, genre.ID)
			if err != nil {
				return fmt.Errorf("failed to delete genres %s", err.Error())
			}

		}
	}
	return nil
}

func (b BookModel) Delete(id int64) error {
	query := `delete from books where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := b.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return err
	}
	return nil
}

func (b BookModel) genresByBook(id int64) ([]*Genre, error) {
	query := `select id, genre_name, created_at, updated_at from genres where id in (select genre_id from books_genres where book_id = $1) order by genre_name`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var genres []*Genre
	rows, err := b.DB.QueryContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	for rows.Next() {
		var genre Genre
		err := rows.Scan(&genre.ID, &genre.GenreName, &genre.CreatedAt, &genre.UpdatedAt)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}
	return genres, nil
}
