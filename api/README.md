## 最小構成

- `cmd/api/main.go` で `.env` を読んで HTTP サーバを起動
- `internal/server/router.go` でパスごとのハンドラを定義
- まだ未実装の API は 501 を返す
- `swagger/openapi.yml` と `swagger/index.html` で Swagger UI / OpenAPI を管理
- OpenAPI は Go の Swagger コメントから自動生成する

## 前提条件

以下のツールがインストールされている必要があります（[Docker Desktop](https://www.docker.com/products/docker-desktop/) を入れると Docker / Docker Compose がまとめて揃います）。

```bash
go version 
make --version         # コマンドのショートカット実行
air -v                 # goのホットリロード
docker -v              # postgresの前提
docker compose version # postgresの前提
```

## 起動

```bash
cp .env.example .env
make db-up
air
```

API はデフォルトで `http://localhost:8080` で起動します。

## Swagger / OpenAPI

- Swagger UI: `http://localhost:8080/swagger/`
- OpenAPI 定義: `swagger/openapi.yml`
- JSON 定義: `swagger/swagger.json`
- 運用メモ: `docs/openapi.md`
