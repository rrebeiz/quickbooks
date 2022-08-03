package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rrebeiz/quickbooks/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrNoRecordFound  = errors.New("resource not found")
	ErrDuplicateEmail = errors.New("duplicate email found")
)

type Users interface {
	GetByEmail(email string) (*User, error)
	GetByID(id int64) (*User, error)
	Insert(user *User) error
	Delete(id int64) error
	Update(user *User) error
	GetAll() ([]*User, error)
	GetAllLoggedIn() ([]*User, error)
}

type Password struct {
	Plaintext *string
	Hash      []byte
}
type User struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    Password  `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	Version     int       `json:"version"`
	AccountType string    `json:"account_type"`
	Token       Token     `json:"token"`
}

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(db *sql.DB) UserModel {
	return UserModel{DB: db}
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Email != "", "email", "should not be empty")
	v.Check(user.Name != "", "name", "should not be empty")
	v.Check(*user.Password.Plaintext != "", "password", "should not be empty")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "should not be empty")
}
func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "should not be empty")
}

func ValidateName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "should not be empty")
}

func (p *Password) HashPassword(plaintext, pepper string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(plaintext+pepper), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintext
	p.Hash = passwordHash
	return nil
}

func (p *Password) CheckPassword(password, pepper string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(password+pepper))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}

	}
	return true, nil
}

func (u UserModel) GetByEmail(email string) (*User, error) {
	query := `select id, name, email, password_hash, created_at, version, account_type from users where email = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user User
	err := u.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.CreatedAt, &user.Version, &user.AccountType)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) Insert(user *User) error {
	query := `insert into users (name, email, password_hash) values($1, $2, $3) returning id, created_at, version`
	args := []interface{}{user.Name, user.Email, user.Password.Hash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (u UserModel) GetAll() ([]*User, error) {
	query := `select id, name, email, password_hash, created_at, version, account_type from users order by name`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var users []*User
	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.CreatedAt, &user.Version, &user.AccountType)
		if err != nil {
			return nil, err
		}
		query := `select id, user_id, email, token, token_hash, created_at, updated_at, expiry from tokens where user_id = $1`
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		var token Token
		err = u.DB.QueryRowContext(ctx, query, user.ID).Scan(&token.ID, &token.UserID, &token.Email, &token.Token, &token.TokenHash, &token.CreatedAt, &token.UpdatedAt, &token.Expiry)
		user.Token = token
		users = append(users, &user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (u UserModel) GetAllLoggedIn() ([]*User, error) {
	query := `select users.id, users.name, users.email, users.password_hash, users.created_at, users.version, users.account_type,
	  tokens.id, tokens.user_id, tokens.email, tokens.token, tokens.token_hash, tokens.created_at, tokens.updated_at, tokens.expiry from users inner join tokens on users.id = tokens.user_id order by name`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []*User
	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.CreatedAt, &user.Version, &user.AccountType, &user.Token.ID,
			&user.Token.UserID, &user.Token.Email, &user.Token.Token, &user.Token.TokenHash, &user.Token.CreatedAt, &user.Token.UpdatedAt, &user.Token.Expiry)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u UserModel) GetByID(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrNoRecordFound
	}
	query := `select id, name, email, password_hash, created_at, version, account_type from users where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user User
	err := u.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.CreatedAt, &user.Version, &user.AccountType)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrNoRecordFound
	}

	query := `delete from users where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := u.DB.ExecContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNoRecordFound
		default:
			return err
		}
	}
	affectedRow, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRow != 1 {
		return ErrNoRecordFound
	}
	return nil
}

func (u UserModel) Update(user *User) error {
	query := `update users set name = $1, email = $2, password_hash = $3, version = version + 1 where id = $4 returning created_at`
	args := []interface{}{user.Name, user.Email, user.Password.Hash, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
