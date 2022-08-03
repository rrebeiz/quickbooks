package data

import "database/sql"

type Models struct {
	Books  Books
	Users  Users
	Tokens Tokens
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books:  NewBookModel(db),
		Users:  NewUserModel(db),
		Tokens: NewTokenModel(db),
	}
}
