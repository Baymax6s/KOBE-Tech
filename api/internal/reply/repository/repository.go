package repository

import (
	"database/sql"
	"errors"
)

var (
	ErrArticleNotFound       = errors.New("article not found")
	ErrParentNotFound        = errors.New("parent reply not found")
	ErrParentMismatch        = errors.New("parent reply does not belong to the article")
	ErrInvalidParent         = errors.New("parent reply kind is not allowed for this reply kind")
	ErrInvalidRootKind       = errors.New("reply kind is not allowed without a parent")
	ErrReplyNotFound         = errors.New("reply not found")
	ErrNotAnswer             = errors.New("reply is not an answer")
	ErrNotQuestionAuthor     = errors.New("only the question author can mark a best answer")
	ErrBestAnswerAlreadySet  = errors.New("best answer already exists for this question")
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
