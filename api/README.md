## 最小構成

- `cmd/api/main.go` で `.env` を読んで HTTP サーバを起動
- `internal/server/router.go` でパスごとのハンドラを定義
- まだ未実装の API は 501 を返す

## 起動

```bash
cp .env.example .env
make db-up
make run
```

API はデフォルトで `http://localhost:8080` で起動します。
