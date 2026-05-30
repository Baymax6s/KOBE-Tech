# Go + Gin のバックエンド API サーバーです。


## 前提条件

以下のツールがインストールされている必要があります（[Docker Desktop](https://www.docker.com/products/docker-desktop/) を入れると Docker / Docker Compose がまとめて揃います）。

```bash
go version
air -v
docker -v
docker compose version
swag --version
```

## 技術スタック

- Go 1.25 / Gin
- PostgreSQL 16
- Swagger (swag) / OpenAPI
- Docker Compose
- air (ホットリロード)


## 命名規則

- ファイル名はsnake_case
- handlerのfile名とrepositoryのfile名は同じ


## 環境構築

```bash
make setup
```

`make setup` で `.env` がなければ `.env.example` から作成します。

## 起動

```sh
make dev
```
または
```sh
go run ./cmd/api/main.go
```

`make dev` で PostgreSQL 起動・マイグレーション・seed データ投入を実行し、`air` で API を起動します。


### DBeaver

| 項目     | 値                |
| -------- | ----------------- |
| Host     | `localhost`       |
| Port     | `5432`            |
| Database | `baymux_db`       |
| Username | `baymux_user`     |
| Password | `baymux_password` |

## 自動生成

Swagger を生成する:

```sh
make swagger
```
または
```sh
swag init -q -g ./cmd/api/main.go -d . --parseInternal -o ./swagger --ot json,yaml
mv ./swagger/swagger.yaml ./swagger/openapi.yml
```

API schema を生成する:

```sh
cd ../app
npm run generate:api
```

## DB関連

DB・マイグレーションを起動 / 停止する:

```sh
make db-up
make db-down
```



マイグレーションファイルを作成する:

```bash
make migrate-create NAME=create_users
```

マイグレーションを手動で巻き戻す:

```bash
make migrate-down
```

### seed データ・ログインアカウント

seed で投入されるデモユーザー一覧、おすすめのログインアカウント、機能ごとの見どころは、リポジトリルートの [README](../README.md#デモデータとログインアカウント) に集約しています（情報の二重管理を避けるため）。デモユーザーのパスワードは全員 `Password` です。


## MinIOの環境設定変更

### MinIO 管理コンソール

URL : http://localhost:9001

Username : minio

Password : password

### MinIO の設定変更
ログイン情報やバケット名など、環境に合わせて設定を変更したい場合は `.env` の以下の項目を書き換えてください。

MinIOのルートアカウント（管理画面のログインに使用）

```env
MINIO_ROOT_USER=minio
MINIO_ROOT_PASSWORD=password
```

アプリケーションからの接続設定（エンドポイント、バケット名、暗号化）

```env
MINIO_ENDPOINT=localhost:9000
MINIO_BUCKET_NAME=kobe-tech-avatar
MINIO_USE_SSL=false
```
