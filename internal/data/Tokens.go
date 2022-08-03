package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"time"
)

type Tokens interface {
	GetByToken(plainText string) (*Token, error)
	GetUserForToken(token *Token) (*User, error)
	GenerateToken(userID int64, ttl time.Duration) (*Token, error)
	InsertToken(token *Token) error
	DeleteToken(id int64) error
}

type Token struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	TokenHash []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Expiry    time.Time `json:"expiry"`
}

type TokenModel struct {
	DB *sql.DB
}

func NewTokenModel(db *sql.DB) TokenModel {
	return TokenModel{DB: db}
}

func (t TokenModel) GetByToken(plainText string) (*Token, error) {
	query := `select id, user_id, email, token, token_hash, created_at, updated_at, expiry from tokens where token = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var token Token
	err := t.DB.QueryRowContext(ctx, query, plainText).Scan(&token.ID, &token.UserID, &token.Email, &token.Token, &token.TokenHash, &token.CreatedAt, &token.UpdatedAt, &token.Expiry)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecordFound
		default:
			return nil, err
		}
	}
	return &token, nil
}

func (t TokenModel) GetUserForToken(token *Token) (*User, error) {
	//query := `select id, name, email, password_hash, created_at, version from users where id = $1`
	query := `select users.id, users.name, users.email, users.password_hash, users.created_at, users.version, users.account_type, tokens.id, tokens.user_id, tokens.email, tokens.token, tokens.token_hash, tokens.created_at, tokens.updated_at, tokens.expiry from users inner join tokens on users.id = tokens.user_id where users.id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := t.DB.QueryRowContext(ctx, query, token.UserID).Scan(&user.ID, &user.Name, &user.Email, &user.Password.Hash, &user.CreatedAt, &user.Version, &user.AccountType, &user.Token.ID, &user.Token.UserID, &user.Token.Email, &user.Token.Token, &user.Token.TokenHash, &user.Token.CreatedAt, &user.Token.UpdatedAt, &user.Token.Expiry)
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

func (t TokenModel) GenerateToken(userID int64, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Token = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Token))
	token.TokenHash = hash[:]

	return token, nil

}

func (t TokenModel) InsertToken(token *Token) error {
	query := `delete from tokens where user_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := t.DB.ExecContext(ctx, query, token.UserID)
	if err != nil {
		return err
	}

	query = `insert into tokens (user_id, email, token, token_hash, expiry) values ($1, $2, $3, $4, $5) returning id, created_at, updated_at`
	args := []interface{}{token.UserID, token.Email, token.Token, token.TokenHash, token.Expiry}
	return t.DB.QueryRowContext(ctx, query, args...).Scan(&token.ID, &token.CreatedAt, &token.UpdatedAt)
}

func (t TokenModel) DeleteToken(id int64) error {
	if id < 1 {
		return ErrNoRecordFound
	}
	query := `delete from tokens where id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)
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

	if affectedRow < 1 {
		return ErrNoRecordFound
	}

	return nil

}
