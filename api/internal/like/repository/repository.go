package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrArticleNotFound = errors.New("article not found")
	ErrAlreadyLiked    = errors.New("already liked")
	ErrUserNotFound    = errors.New("user not found")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
