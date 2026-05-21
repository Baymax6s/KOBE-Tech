# 環境構築

## 前提条件

以下のツールがインストールされている必要があります（[Docker Desktop](https://www.docker.com/products/docker-desktop/) を入れると Docker / Docker Compose がまとめて揃います）。

```bash
go version
air -v
docker -v
docker compose version
swag --version
```

## getting started

```bash
make setup
make dev
```

`make setup` で `.env` がなければ `.env.example` から作成します。
`make dev` で PostgreSQL 起動・マイグレーション・seed データ投入を実行し、`air` で API を起動します。

## 起動

```sh
make dev
```

## DB確認

`lazysql`または`DBeaver`でDBを確認することを推奨します

### lazysql

`./scripts/lazysql.sh` で `DATABASE_URL` に接続した `lazysql` を read-only で開けます。

オプションを渡す場合は `./scripts/lazysql.sh --help` のように実行してください。`lazysql` 本体は各開発環境にインストールされている必要があります。

### DBeaver

| 項目     | 値                |
| -------- | ----------------- |
| Host     | `localhost`       |
| Port     | `5432`            |
| Database | `baymux_db`       |
| Username | `baymux_user`     |
| Password | `baymux_password` |

## コマンド一覧

Swagger を生成する:

```bash
make swagger
```

DB・マイグレーションを起動 / 停止する:

```bash
make db-up
make db-down
```

テストを実行する:

```bash
make test
```

マイグレーションを手動で巻き戻す:

```bash
make migrate-down
```

マイグレーションファイルを作成する:

```bash
make migrate-create NAME=create_users
```

### seed_users_データ

| user_id | name   | password |
| ------- | ------ | -------- |
| 1       | admin  | Password |
| 2       | user01 | Password |
| 3       | user02 | Password |
| 4       | user03 | Password |
