# KOBE Tech Frontend

Vue 3 + Vite のフロントエンドです。`src` は `page / components / api` を中心にした簡素な構成です。

## Setup

```sh
npm install
```

## Development

```sh
npm run dev
```

## Build

```sh
npm run build
```

## Directory Layout

```text
src/
  main.ts
  main.css      # 全体スタイル
  App.vue
  page/         # 画面単位
  components/   # 画面配下の表示部品
  api/          # データ取得
```

## Dependency Rule

- `page` は `components` と `api` を読む
- `components` は props だけ受け取る
- `api` は取得処理と型を持つ
