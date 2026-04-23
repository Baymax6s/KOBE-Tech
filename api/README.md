## 最小構成

- `cmd/api/main.go` で `.env` を読んで HTTP サーバを起動
- `internal/server/router.go` でパスごとのハンドラを定義
- まだ未実装の API は 501 を返す
- `swagger/openapi.yml` と `swagger/index.html` で Swagger UI / OpenAPI を管理
- OpenAPI は Go の Swagger コメントから自動生成する

## 起動

```bash
cp .env.example .env
make db-up
make run
```

API はデフォルトで `http://localhost:8080` で起動します。

## Swagger / OpenAPI

- Swagger UI: `http://localhost:8080/swagger/`
- OpenAPI 定義: `swagger/openapi.yml`
- JSON 定義: `swagger/swagger.json`
- 運用メモ: `docs/openapi.md`
