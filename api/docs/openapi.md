# OpenAPI の運用

- OpenAPI 定義は Go の Swagger コメントから自動生成する
- 生成先は `swagger/openapi.yml` と `swagger/swagger.json`
- Swagger UI は `swagger/index.html` で表示する
- API 実装を追加または変更したら、同じタイミングで Swagger コメントも更新する

## 確認方法

```bash
make swagger
make run
```

ブラウザで `http://localhost:8080/swagger/` を開くと Swagger UI を確認できます。

## PR コメント

- `.github/workflows/preview.yml` が `api/swagger` を GitHub Pages preview に公開し、PR コメントを毎回作り直します
- `.github/workflows/cleanup.yml` が PR close 時に preview を削除します
- 構成は `~/github/optical-backend` の preview-pages ベースに合わせています

## ローカル生成

```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.4
make swagger
```

`make build` と `make run` は自動で `make swagger` を先に実行します。
