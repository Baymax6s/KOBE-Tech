# KOBE Tech Frontend

Vue 3 + Vite のフロントエンドです。`src` は `features / components` を中心にした簡素な構成です。
apiSchemaは自動生成されます。

## Setup

### versionマネージャーの導入推奨
例 nvm mise asdfなど..

node.js versionは24
```sh
node --version
```

### 初期setup

```sh
npm install
```

### frontendの起動コマンド

```sh
npm run dev
```

### Generate API Types

OpenAPIファイルから、TypeScriptクライアントを生成します。

```sh
npm run generate:api
```
