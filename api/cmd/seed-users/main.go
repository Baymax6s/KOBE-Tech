package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	Name     string
	Password string
}

const defaultSeedPassword = "Password"

var requiredUsers = []seedUser{
	{
		Name:     "admin",
		Password: defaultSeedPassword,
	},
	{
		Name:     "user01",
		Password: defaultSeedPassword,
	},
	{
		Name:     "user02",
		Password: defaultSeedPassword,
	},
	{
		Name:     "user03",
		Password: defaultSeedPassword,
	},
}

const insertUserQuery = `
INSERT INTO users (name, password_hash)
VALUES ($1, $2)
`

const updateUserPasswordQuery = `
UPDATE users
SET password_hash = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE name = $2
`

func main() {
	if err := loadEnvFile(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("load .env: %v", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	for _, user := range requiredUsers {
		action, err := seedRequiredUser(ctx, db, user)
		if err != nil {
			log.Fatalf("seed user %q: %v", user.Name, err)
		}

		log.Printf("user=%s action=%s", user.Name, action)
	}
}

func seedRequiredUser(ctx context.Context, db *sql.DB, user seedUser) (string, error) {
	const selectUserQuery = `
SELECT password_hash
FROM users
WHERE name = $1
`

	var currentHash string
	err := db.QueryRowContext(ctx, selectUserQuery, user.Name).Scan(&currentHash)
	if errors.Is(err, sql.ErrNoRows) {
		hash, hashErr := hashPassword(user.Password)
		if hashErr != nil {
			return "", hashErr
		}

		if _, execErr := db.ExecContext(ctx, insertUserQuery, user.Name, hash); execErr != nil {
			return "", execErr
		}

		return "inserted", nil
	}

	if err != nil {
		return "", err
	}

	matches, err := passwordHashMatches(currentHash, user.Password)
	if err != nil {
		return "", fmt.Errorf("compare password hash: %w", err)
	}

	if matches {
		return "unchanged", nil
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}

	if _, err := db.ExecContext(ctx, updateUserPasswordQuery, hash, user.Name); err != nil {
		return "", err
	}

	return "updated", nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate password hash: %w", err)
	}

	return string(hash), nil
}

func passwordHashMatches(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, err
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(strings.TrimPrefix(line, "export "), "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}

		if len(value) >= 2 {
			if value[0] == '"' && value[len(value)-1] == '"' {
				value = value[1 : len(value)-1]
			}
			if value[0] == '\'' && value[len(value)-1] == '\'' {
				value = value[1 : len(value)-1]
			}
		}

		if _, exists := os.LookupEnv(key); exists {
			continue
		}

		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}
