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

### コード品質コマンド

```sh
npm run lint
npm run format
npm run fix
```

`lint` は `type-check` と ESLint を実行します。`format` は Prettier で書き込み、`fix` は format と ESLint の自動修正までまとめて実行します。

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

### MSW の有効化

MSW(Mock Service Worker) を利用する場合は `.env` を編集してください。

```env
VITE_USE_MSW=true
```

デフォルトでは `false` になっており、実APIを利用します。

MSW を利用することで、backend 未実装時でも frontend 単体で動作確認できます。

frontend の実装では、API 実装に加えて MSW のモック追加までをタスク範囲とします。

## MSW 開発

### ファイル構成

```txt
app/
├── .env                  # VITE_USE_MSW=true / false
└── src/
    ├── mocks/
    │   ├── handlers.ts   # MSW handler / モックAPI / モックデータ
    │   └── browser.ts    # MSW worker セットアップ
    └── main.ts           # .env を利用した MSW 有効化設定
```

### モックAPI作成

モックAPIは基本的に `src/mocks/handlers.ts` で管理してください。

`handlers.ts` では以下をまとめて管理しています。

- MSW handler
- モックAPI
- モックデータ
- utility関数

handler は以下の形式で定義してください。

```ts
http.get('*/api/articles', ...)
http.post('*/api/auth/login', ...)
```

frontend の実装では、API 実装に加えて MSW のモック追加までをタスク範囲とします。
