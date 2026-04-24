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
)

type seedArticle struct {
	Title    string
	Content  string
	UserName string
}

var seedArticles = []seedArticle{
	{
		Title:    "神戸大学でのハッカソン体験記",
		Content:  "先日、神戸大学で開催されたハッカソンに参加しました。チームメンバーと48時間かけてWebアプリを開発し、多くのことを学びました。特にチーム開発の難しさと楽しさを実感できた貴重な経験でした。",
		UserName: "admin",
	},
	{
		Title:    "Goで作るREST API入門",
		Content:  "GoのginフレームワークでREST APIを作る基本を紹介します。ルーティング、ミドルウェア、JSONレスポンスの返し方など、実際のコードを交えながら解説します。Goは静的型付けと高いパフォーマンスが魅力で、バックエンド開発に最適です。",
		UserName: "user01",
	},
	{
		Title:    "Vue 3 + Vuetifyで学ぶフロントエンド開発",
		Content:  "Vue 3のComposition APIとVuetifyを組み合わせたフロントエンド開発の入門記事です。コンポーネント設計やリアクティブな状態管理の基本を、サンプルコードを通じて説明します。",
		UserName: "user02",
	},
	{
		Title:    "PostgreSQLのマイグレーション管理",
		Content:  "golang-migrateを使ったDBマイグレーションの管理方法を解説します。up/downファイルの書き方やCIへの組み込み方など、チーム開発で役立つプラクティスを紹介します。",
		UserName: "user03",
	},
	{
		Title:    "Dockerで開発環境を統一する",
		Content:  "docker composeを使って開発環境を統一する方法を紹介します。PostgreSQLのコンテナ起動やボリューム管理など、チームで環境差異をなくすためのTipsをまとめました。",
		UserName: "user01",
	},
}

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

	for _, article := range seedArticles {
		action, err := upsertArticle(ctx, db, article)
		if err != nil {
			log.Fatalf("seed article %q: %v", article.Title, err)
		}
		log.Printf("title=%q user=%s action=%s", article.Title, article.UserName, action)
	}
}

func upsertArticle(ctx context.Context, db *sql.DB, article seedArticle) (string, error) {
	const selectUserQuery = `SELECT id FROM users WHERE name = $1`
	var userID int64
	if err := db.QueryRowContext(ctx, selectUserQuery, article.UserName).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user %q not found: run seed-users first", article.UserName)
		}
		return "", err
	}

	const existsQuery = `SELECT id FROM articles WHERE title = $1 AND user_id = $2`
	var existingID int64
	err := db.QueryRowContext(ctx, existsQuery, article.Title, userID).Scan(&existingID)
	if err == nil {
		return "skipped", nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	const insertQuery = `
		INSERT INTO articles (title, content, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`
	if _, err := db.ExecContext(ctx, insertQuery, article.Title, article.Content, userID); err != nil {
		return "", err
	}

	return "inserted", nil
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
