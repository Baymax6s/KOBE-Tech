package postlogin

import "time"

type User struct {
	ID           int64
	Name         string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}