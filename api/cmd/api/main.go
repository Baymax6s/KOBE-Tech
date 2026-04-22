package main

import (
	"bufio"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Baymax6s/KOBE-Tech/api/internal/server"
)

// @title KOBE-Tech API
// @version 0.1.0
// @description KOBE-Tech API documentation generated from Go annotations.
// @BasePath /
func main() {
	if err := loadEnvFile(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("load .env: %v", err)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("api listening on %s", addr)

	log.Fatal(http.ListenAndServe(addr, server.NewHandler()))
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
