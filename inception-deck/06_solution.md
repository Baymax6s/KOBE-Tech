# 06. 解決案を描く

> 状態: 確定（現リポジトリの構成に基づく）

技術的な全体像をざっくり共有する。詳細仕様ではなく「どんな部品でできているか」を握るためのもの。

## 全体構成

```
[ブラウザ / Vue 3 SPA]
        │  HTTP (axios) / 自動生成APIクライアント
        ▼
[Go + Gin API]
        ├── Postgres        … users / articles / tags / likes / replies / profile
        └── MinIO            … 画像オブジェクト（アバター・添付）
```

## フロントエンド

- Vue 3.5 + TypeScript + Vite
- UI: Vuetify 4
- 状態管理: Pinia 3 / ルーティング: Vue Router 4
- Markdown: `markdown-it` + `md-editor-v3` + `highlight.js`（記事は Markdown）
- 画像トリミング: `vue-advanced-cropper`
- 機能ディレクトリ: `articles / auth / profile / replies / settings / home / about`

## バックエンド

- Go 1.25 + Gin
- 認証: JWT（`golang-jwt/jwt`）+ bcrypt（`golang.org/x/crypto`）
- DB: Postgres（`lib/pq`）
- オブジェクトストレージ: MinIO（`minio-go`、画像の presign）
- ドメイン: `article / auth / like / profile / reply / minio / server`

## コントラクト層（重要）

API の型は `api/swagger/openapi.yml` を真のソースとし、`app/src/api/generated/` へ自動生成する。

```
Go handler/DTO を編集
  → cd api && make swagger      （OpenAPI 更新）
  → cd app && npm run generate:api （TSクライアント再生成）
```

`app/src/api/generated/**` は手で編集しない。

## 中核（差別化）の実装位置

差別化の要「透明化＝プロフィールの充実」は、**既存の `profile` ドメイン / `profile` 機能を厚くする**ことで実現する。
新しい大型機能を足すのではなく、bio・制作物・投稿履歴をプロフィールに集約し、投稿者名リンク（既存導線）から辿れる状態にする。
