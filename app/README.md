# KOBE Tech Frontend

Vue 3 + Vite のフロントエンドです。`src` は `features / components` を中心にした簡素な構成です。
apiSchemaは`npm run generate:api`によって生成されます。

## Setup

### versionマネージャーの導入推奨

例 nvm mise asdfなど..(推奨はmise)

node.js ltsを採用 v24.15.0

```sh
node --version
```

### 初期setup

```sh
cp .env.example .env
npm install
```

### frontendの起動コマンド

```sh
npm run dev
```

### Generate API Types

`api/swagger/openapi.yml` から TypeScript クライアントを再生成します。
backend 側の OpenAPI 定義を更新した直後であれば、先に以下を実行してください。

```sh
cd ../api
swag init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
mv ./swagger/swagger.yaml ./swagger/openapi.yml
```

```sh
npm run generate:api
```
