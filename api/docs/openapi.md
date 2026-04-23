# OpenAPI の運用

- OpenAPI 定義は Go の Swagger コメントから自動生成する
- 生成先は `swagger/openapi.yml` と `swagger/swagger.json`
- Swagger UI は `swagger/index.html` で表示する
- API 実装を追加または変更したら、同じタイミングで Swagger コメントも更新する

## 確認方法

```bash
swag init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
mv ./swagger/swagger.yaml ./swagger/openapi.yml
go run ./cmd/api
```

ブラウザで `http://localhost:8080/swagger/` を開くと Swagger UI を確認できます。

## PR コメント

- `.github/workflows/preview.yml` が `api/swagger` を GitHub Pages preview に公開し、PR コメントを毎回作り直します
- `.github/workflows/cleanup.yml` が PR close 時に preview を削除します
- preview 公開と PR close 時の cleanup は、このリポジトリ内の workflow 設定に沿って運用しています

## ローカル生成

```bash
go install github.com/swaggo/swag/cmd/swag@v1.16.4
swag init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
mv ./swagger/swagger.yaml ./swagger/openapi.yml
```

`go build` と `go run` は Swagger を自動生成しないため、Swagger コメントを更新したときは先に再生成してください。
