package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mozillazg/go-slugify"
	"github.com/rrebeiz/quickbooks/internal/validator"
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
	Reviews         []*Review `json:"reviews,omitempty"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}
type Author struct {
	ID         int64     `json:"id"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
	Version    int       `json:"-"`
}

type Genre struct {
	ID        int64     `json:"id"`
	GenreName string    `json:"genre_name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
type Review struct {
	ID        int64     `json:"id"`
	Rating    int       `json:"rating"`
	Review    string    `json:"review"`
	BookID    int64     `json:"-"`
	Book      string    `json:"book,omitempty"`
	UserID    int64     `json:"-"`
	User      string    `json:"user,omitempty"`
	Version   int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Books interface {
	GetAll(title string, filters Filters) ([]*Book, Metadata, error)
	GetByID(id int64) (*Book, error)
	GetBySlug(slug string) (*Book, error)
	Insert(book *Book) error
	Update(book *Book) error
	Delete(id int64) error
	GetAllAuthors(author string, filters Filters) ([]*Author, Metadata, error)
	GetAllReviewsByUser(user string, filters Filters) ([]*Review, Metadata, error)
	GetReviewByID(id int64) (*Review, error)
	InsertReview(review *Review) error
	UpdateReview(review *Review) error
	DeleteReview(id int64) error
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

func ValidateReview(v *validator.Validator, review *Review) {
	v.Check(review.Rating > 0, "rating", "should not be empty")
	v.Check(review.Rating <= 5, "rating", "should not be bigger than 5")
	v.Check(review.Review != "", "review", "should not be empty")
}

func (b BookModel) GetAll(title string, filters Filters) ([]*Book, Metadata, error) {
	query := fmt.Sprintf(`select count (*) over(), b.id, b.title, b.author_id, b.publication_year, b.slug, b.description, b.created_at, 
						b.updated_at, a.id, a.author_name, a.created_at, a.updated_at, a.version from books b left join authors a on (b.author_id = a.id) 
						where (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')  
						order by b.%s %s, b.id ASC limit $2 offset $3`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var books []*Book
	args := []interface{}{title, filters.limit(), filters.offset()}

	rows, err := b.DB.QueryContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, Metadata{}, ErrNoRecordFound
		default:
			return nil, Metadata{}, err
		}
	}
	totalRecords := 0

	for rows.Next() {
		var book Book
		err := rows.Scan(&totalRecords, &book.ID, &book.Title, &book.AuthorID, &book.PublicationYear, &book.Slug, &book.Description, &book.CreatedAt,
			&book.UpdatedAt, &book.Author.ID, &book.Author.AuthorName, &book.Author.CreatedAt, &book.Author.UpdatedAt, &book.Author.Version)
		if err != nil {
			return nil, Metadata{}, err
		}

		genres, err := b.genresByBook(book.ID)
		if err != nil {
			return nil, Metadata{}, err
		}
		reviews, err := b.reviewsByBook(book.ID)
		if err != nil {
			return nil, Metadata{}, err
		}
		for i := range genres {
			book.Genres = append(book.Genres, genres[i].GenreName)
		}
		book.Reviews = reviews
		books = append(books, &book)
	}
	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return books, metadata, nil
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
	reviews, err := b.reviewsByBook(book.ID)
	if err != nil {
		return nil, err
	}
	var bookGenres []string
	for _, x := range genres {
		bookGenres = append(bookGenres, x.GenreName)
	}
	book.Genres = bookGenres
	book.Reviews = reviews
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

func (b BookModel) GetAllAuthors(author string, filters Filters) ([]*Author, Metadata, error) {
	query := fmt.Sprintf(`select count(*) over(), id, author_name from authors where (to_tsvector('simple', author_name) @@ plainto_tsquery('simple', $1) OR $1 = '') 
			order by %s %s, id asc limit $2 offset $3`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var authors []*Author
	args := []interface{}{author, filters.limit(), filters.offset()}
	rows, err := b.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	for rows.Next() {
		var author Author
		err := rows.Scan(&totalRecords, &author.ID, &author.AuthorName)
		if err != nil {
			return nil, Metadata{}, err
		}
		authors = append(authors, &author)
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return authors, metadata, nil

}
func (b BookModel) GetAllReviewsByUser(user string, filters Filters) ([]*Review, Metadata, error) {
	query := fmt.Sprintf(`select count(*) over(), u.id, u.name, b.id, b.title, b.author_id, b.publication_year, r.id, r.rating, r.review, r.user_id, r.book_id from users u join reviews r on u.id = r.user_id join books b on b.id = r.book_id where (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') order by u.%s %s, u.id asc limit $2 offset $3`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var reviews []*Review
	args := []interface{}{user, filters.limit(), filters.offset()}
	rows, err := b.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	for rows.Next() {
		var review Review
		var user User
		var book Book
		err := rows.Scan(&totalRecords, &user.ID, &user.Name, &book.ID, &book.Title, &book.AuthorID, &book.PublicationYear, &review.ID, &review.Rating, &review.Review, &review.UserID, &review.BookID)
		if err != nil {
			return nil, Metadata{}, err
		}
		review.Book = book.Title
		review.User = user.Name
		reviews = append(reviews, &review)

	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return reviews, metadata, nil
}

func (b BookModel) GetReviewByID(id int64) (*Review, error) {
	query := `select id, rating, review, book_id, user_id, version, created_at, updated_at from reviews where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var review Review
	err := b.DB.QueryRowContext(ctx, query, id).Scan(&review.ID, &review.Rating, &review.Review, &review.BookID, &review.UserID, &review.Version, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	return &review, nil
}

func (b BookModel) InsertReview(review *Review) error {
	query := `insert into reviews (rating, review, book_id, user_id) values ($1, $2, $3, $4) returning id, version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{review.Rating, review.Review, review.BookID, review.UserID}
	return b.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID, &review.Version)
}

func (b BookModel) UpdateReview(review *Review) error {
	query := `update reviews set rating = $1, review = $2, updated_at = now(), version = version + 1 where id = $3 and version = $4 returning version`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := b.DB.QueryRowContext(ctx, query, review.Rating, review.Review, review.ID, review.Version).Scan(&review.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNoRecordFound
		default:
			return err
		}
	}
	return nil
}

func (b BookModel) DeleteReview(id int64) error {
	query := `delete from reviews where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := b.DB.ExecContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNoRecordFound
		default:
			return err
		}
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row != 1 {
		return ErrNoRecordFound
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

func (b BookModel) reviewsByBook(id int64) ([]*Review, error) {
	query := `select u.id, u.name, r.id, r.rating, r.review, r.book_id, r.user_id, 
       r.version, r.created_at, r.updated_at from users u join reviews r on u.id = r.user_id where book_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var reviews []*Review

	rows, err := b.DB.QueryContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	defer rows.Close()
	for rows.Next() {
		var review Review
		var user User
		err := rows.Scan(&user.ID, &user.Name, &review.ID, &review.Rating, &review.Review, &review.BookID, &review.UserID, &review.Version, &review.CreatedAt, &review.UpdatedAt)
		if err != nil {
			return nil, err
		}
		review.User = user.Name
		reviews = append(reviews, &review)
	}
	return reviews, nil
}
