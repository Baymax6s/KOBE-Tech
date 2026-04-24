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
docker compose up -d
swag init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
mv ./swagger/swagger.yaml ./swagger/openapi.yml
air
```

`docker compose up -d` で PostgreSQL 起動・マイグレーション・seed データ投入がまとめて実行されます。

## 起動
```sh
docker compose up -d
air
```

## DB確認

`./scripts/lazysql.sh` で `DATABASE_URL` に接続した `lazysql` を read-only で開けます。

オプションを渡す場合は `./scripts/lazysql.sh --help` のように実行してください。`lazysql` 本体は各開発環境にインストールされている必要があります。

## コマンド一覧

Swagger を生成する:

```bash
swag init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
mv ./swagger/swagger.yaml ./swagger/openapi.yml
```

DB・マイグレーションを起動 / 停止する:

```bash
docker compose up -d
docker compose down
```

マイグレーションを手動で巻き戻す:

```bash
docker compose run --rm migrate down 1
```

マイグレーションファイルを作成する:

```bash
docker compose run --rm migrate create -ext sql -dir /migrations <ファイル名のサフィックス>
```
