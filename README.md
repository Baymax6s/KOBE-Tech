## 環境構築

```sh
cp .env.example .env
docker compose up -d postgres
go install github.com/swaggo/swag/cmd/swag@v1.16.4
swag init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
mv ./swagger/swagger.yaml ./swagger/openapi.yml
go run ./cmd/api
```

## 起動コマンド
```sh
docker compose up -d
air
```

## 終了コマンド
```sh
Ctrl + C
docker compose down
```
