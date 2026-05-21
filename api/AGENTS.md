# api/ ルール（Go バックエンド）

`api/` 配下を編集するすべての AI エージェントはここを参照してください。

## エントリポイントとフレームワーク

- エントリポイント: `cmd/api/main.go`（デフォルトポート 8080）
- Gin 1.12, Go 1.25, `github.com/Baymax6s/KOBE-Tech/api`
- ホットリロード: `air`（`.air.toml`）

## パッケージ構造

```
internal/
  article/, auth/, like/, profile/, reply/  → handler/ + repository/
  auth/    → さらに jwt/issuer/validator/password
  server/  → router.go / cors.go / swagger.go
```

- 機能単位で `handler/`（リクエスト処理）と `repository/`（DB アクセス）に分割
- `server/` が依存性注入 + ルーティング
- ファイル名は snake_case、**handler と repository で同じファイル名**（例: `handler/get_article.go` / `repository/get_article.go`）

## ハンドラの規約

```go
type Handler struct { repo *repository.Repository }
func NewHandler(repo *repository.Repository) *Handler { ... }
func (h *Handler) RegisterRoutes(router gin.IRouter, authRouter gin.IRouter) { ... }
```

- 公開 → `router`、認証必須 → `authRouter`
- swag アノテーションを各 handler メソッドに付与
- DTO は handler ファイル内で定義、`// @name server.xxx` で型名を明示

## 認証ミドルウェア

- `auth.RequireUser(v)` → `auth.MustUserID(c)` で userID 取得
- `auth.OptionalUser(v)` → `auth.OptionalUserID(c)` で取得
- JWT: HMAC-SHA256, 24h 期限, claims に `user_id`
- パスワード: bcrypt

## DB・Swagger

- PostgreSQL 16 + `database/sql` + `lib/pq`。コネクションプール: MaxOpenConns=10, Idle=10, Lifetime=5min
- マイグレーション: `migrate/migrate` v4（Docker Compose 内で自動実行）
- Swagger: `//go:embed` で埋め込み、`/swagger/` で配信
- `make swagger` → `swag init` + `swagger.yaml` を `openapi.yml` にリネーム
- `swagger.yaml` は .gitignore、`openapi.yml` / `swagger.json` はコミット対象

## 主要コマンド

| コマンド | 用途 |
|----------|------|
| `make setup` | `.env.example` → `.env` |
| `make dev` | PostgreSQL 起動 + air |
| `make swagger` | OpenAPI 再生成 |
| `make db-up/down` | DB 起動 / 停止 |
| `make migrate-create NAME=x` | マイグレーションファイル作成 |
| `make migrate-down` | 1 つ巻き戻し |
| `go vet ./...` | 静的解析 |
| `go build -o bin/api ./cmd/api` | ビルド |

CI では `gofmt -l` + `go vet ./...`。テストは未整備。

## CORS・環境変数

- デフォルトで localhost, `baymax6s.github.io`, `vue-cjne.onrender.com` を許可
- `CORS_ALLOWED_ORIGINS` で追加（カンマ区切り）。`CORS_MAX_AGE_SECONDS` でキャッシュ制御
- `.env` を読み込む（.gitignore）。**`JWT_SECRET` は必須**

| 変数 | デフォルト | 必須 |
|------|-----------|------|
| `DATABASE_URL` | — | yes |
| `JWT_SECRET` | — | yes |
| `APP_PORT` | 8080 | no |
| `CORS_ALLOWED_ORIGINS` | — | no |
| `CORS_MAX_AGE_SECONDS` | 0 | no |

## 注意点

- `swagger/swagger.yaml` は .gitignore（生成中間物）。`bin/`, `.codex` も gitignore
- **`swagger/openapi.yml` / `swagger/swagger.json` は生成ファイル。ハンドラの swag アノテーションを編集し `make swagger` で再生成すること。直接編集禁止**
- `swag` のインストール: `go install github.com/swaggo/swag/cmd/swag@v1.16.4`
- Docker Compose プロジェクト名: `baymux`
- 編集後の自動フォーマット無し（Claude / OpenCode hooks 任せ）
