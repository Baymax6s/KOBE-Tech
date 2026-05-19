package auth

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidLoginRequest = errors.New("name and password are required")

type User struct {
	ID           int64
	Name         string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NormalizeLoginInput(name, password string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || password == "" {
		return "", ErrInvalidLoginRequest
	}

	return name, nil
}
