package profile

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
	"database/sql"
)

type User struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserProfile struct {
	ID         int64
	UserID     int64
	ObjectKey  sql.NullString
	Bio        sql.NullString
	IsUploaded bool
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
}

type Profile struct {
	User        User
	UserProfile UserProfile
}

const maxBioLength = 200

var (
	ErrInvalidBio = errors.New("bio is required")
	ErrBioTooLong = errors.New("bio must be 200 characters or less")
)

func NormalizeBioInput(bio string) (string, error) {
	bio = strings.TrimSpace(bio)

	if utf8.RuneCountInString(bio) > maxBioLength {
		return "", ErrBioTooLong
	}

	return bio, nil
}
