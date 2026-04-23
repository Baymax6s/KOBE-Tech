## 前提条件

以下のツールがインストールされている必要があります（[Docker Desktop](https://www.docker.com/products/docker-desktop/) を入れると Docker / Docker Compose がまとめて揃います）。

```bash
go version
air -v
docker -v
docker compose version
swag --version
```

## 環境構築

```bash
cp .env.example .env
docker compose up -d postgres
swag init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
mv ./swagger/swagger.yaml ./swagger/openapi.yml
air
```

## 起動
```sh
docker compose up -d
air
```

## ユーザーseed

`go run ./cmd/seed-users` はログイン用ユーザーを投入または更新します。デフォルトでは `admin`, `user01`, `user02`, `user03` を同期し、初期パスワードはすべて `Password` です。パスワードは seed 実行時に bcrypt でハッシュ化して保存します。

## コマンド一覧

Swagger を生成する:

```bash
swag init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
mv ./swagger/swagger.yaml ./swagger/openapi.yml
```

DB を起動 / 停止する:

```bash
docker compose up -d postgres
docker compose down
```

マイグレーションを適用 / 巻き戻しする:

```bash
docker compose --profile tools run --rm migrate up
docker compose --profile tools run --rm migrate down 1
```

マイグレーションファイルを作成する:

```bash
docker compose --profile tools run --rm migrate create -ext sql -dir /migrations -seq create_users
```
