# プロフィール画面リニューアル 要件定義（ドラフト / 壁打ち用）

> ✅ 論点は確定済み（「確定事項」セクション参照）。この内容で実装に入れる状態。

## 概要

他人のプロフィールが閲覧できるようになった（記事詳細の著者名リンク経由）。
これに合わせてプロフィール画面のレイアウトを刷新し、1 ユーザーの活動が一望できる画面にする。

## ゴール

- プロフィール画面を「ほぼ全画面」のレイアウトに刷新する
  - **上部（ヘッダー）**: アバターアイコン・名前・自己紹介
  - **下部（タブ）**: 「投稿した記事」一覧 / 「いいねした記事」一覧の 2 タブ
- 自分のプロフィール（`/profile/me`）と他人のプロフィール（`/profile/:userId`）を同じ画面で扱う
- 自分の場合のみ、自己紹介の編集（とアバター変更）ができる
- バックエンドに「指定ユーザーの投稿記事一覧」「指定ユーザーのいいね記事一覧」を返す API を追加する

## ノンゴール（今回やらない）

- 記事一覧のページネーション（初期は全件一括取得。記事一覧画面も現状ページネーションなし）
- フォロー / フォロワー
- プロフィールの公開・非公開設定
- いいねの取り消しをこの画面から行う動線
- 記事の編集・削除動線

---

## 現状の把握（コードベース調査結果）

実装方針を決める前提として、今あるものを整理する。

### フロント

- `app/src/features/profile/ProfileView.vue`
  - `/profile/me`（`isMe` prop）と `/profile/:userId`（`userId` prop）の両方を担う 1 コンポーネント
  - アバターは `mdi-account-circle` 固定アイコン（**画像表示は未実装**）
  - `is_owner` のときだけ自己紹介の編集ボタンを表示
- `app/src/features/articles/ArticleCard.vue`
  - タイトル・タグ・投稿日・いいね数を表示する再利用可能なカード。**著者名は表示していない**
- `app/src/features/articles/ArticlesView.vue`
  - 絞り込み状態を `route.query.tag` に持たせる設計（リロード・共有で再現可能）。タブ状態の設計もこれに倣える
- ルーティングは `/profile/me`（`requiresAuth`）と `/profile/:userId(\d+)`（公開）

### バックエンド

- `GET /api/profile/{user_id}` … `id` / `name` / `bio` / `is_owner` を返す。**アバターURL は返していない**
  - `user_profiles` テーブルに `object_key`・`is_uploaded` はあるが、取得 API では使っていない
  - アバターのアップロードは `POST /api/profile/avatar/presign` → `/complete` で実装済み（MinIO）
- `GET /api/articles` … タグ絞り込み（AND）のみ。**`user_id` での絞り込みは未対応**
  - 一覧でも本文 `content` を**全文**返している（プロフィールで記事が増えると転送量が無駄）
- いいねは `likes(article_id, user_id, created_at)`。**「ユーザーのいいね記事一覧」を返す API は未実装**

### スキーマ（関連分）

```
users(id, name, password_hash, ...)
user_profiles(id, user_id UNIQUE, object_key, bio, is_uploaded, created_at, updated_at)
articles(id, title, content, user_id, created_at, updated_at)
likes(id, article_id, user_id, created_at, UNIQUE(article_id, user_id))
```

---

## 画面要件（ドラフト）

```text
┌─────────────────────────────────────┐
│  [アバター]  名前                      │  ← ヘッダー
│              自己紹介テキスト…          │     (自分なら「編集」ボタン)
├─────────────────────────────────────┤
│  [ 投稿した記事 ]  [ いいねした記事 ]    │  ← v-tabs
├─────────────────────────────────────┤
│  ArticleCard                         │  ← v-window で切替
│  ArticleCard                         │
│  ...                                 │
└─────────────────────────────────────┘
```

- レイアウトは Vuetify のみで構成（`v-tabs` + `v-window` + `v-card`）。`min-h-screen` は使わず `v-container fill-height` 系で縦を確保（app/AGENTS.md ルール）
- 記事カードは既存 `ArticleCard.vue` を流用
- 各タブの空状態（投稿 0 件 / いいね 0 件）を `v-alert` で出す
- 投稿/いいねの一覧はプロフィール取得後に並行取得し、ローディング・エラーは 2 タブ共通で扱う（同じネットワーク要因で同時に成否が決まることがほとんどなので、タブごとに loading/error フラグを増やさず宣言的に保つ）

---

## バックエンド設計（ドラフト）

### 追加エンドポイント案

| メソッド | パス | 用途 |
| -------- | ---- | ---- |
| `GET` | `/api/profile/{user_id}/articles` | そのユーザーが投稿した記事一覧 |
| `GET` | `/api/profile/{user_id}/liked-articles` | そのユーザーがいいねした記事一覧 |

- レスポンスは既存の `listArticlesResponse`（`ArticleListItemJSON`）を再利用する想定
- `liked_by_me` は「閲覧者(current user)が」いいね済みか。プロフィール対象ユーザーのいいねとは別概念なので混同しない

### クエリの骨子

- 投稿一覧: `WHERE a.user_id = $target`
- いいね一覧: `WHERE a.id IN (SELECT article_id FROM likes WHERE user_id = $target)`
  - 並び順は「いいねした順（`likes.created_at DESC`）」か「記事投稿順」かが論点（下記）

---

## 確定事項

| # | 論点 | 決定 |
| - | ---- | ---- |
| 1 | 投稿記事一覧 API の置き場所 | **profile 配下に新規** `GET /api/profile/{user_id}/articles` |
| 2 | 他人の「いいね一覧」公開 | **全員に公開**（`/api/profile/{user_id}/liked-articles`） |
| 3 | アバター画像表示 | **今回スコープに含める**。プロフィール取得 API に `avatar_url` を追加 |
| 4 | タブ状態 | **URL に保持**（`?tab=articles\|likes`） |
| 5 | いいね一覧の並び順 | **いいねした順**（`likes.created_at DESC`） |
| 6 | 一覧 API が本文を全文返す件 | **現状踏襲**（スコープを膨らませない。別途リファクタ扱い） |

### 論点 3 の詳細（アバター）

- `GET /api/profile/{user_id}` のレスポンスに `avatar_url` を追加する
- `is_uploaded = true` のとき、`object_key` から **presigned GET URL** を生成して返す（既存 MinIO クライアントを利用）
- 未アップロード時は `avatar_url` を空にし、フロントは従来どおり `mdi-account-circle` にフォールバック

### 論点 1・2 の詳細（一覧 API）

- どちらも既存の `listArticlesResponse`（`ArticleListItemJSON`）を再利用
- `liked_by_me` は閲覧者(current user)基準。プロフィール対象ユーザーのいいねとは別概念
- 並び順
  - 投稿一覧: `articles.created_at DESC`
  - いいね一覧: `likes.created_at DESC`

---

## 想定する作業の流れ（確定後）

1. マイグレーション不要（既存テーブルで足りる想定）
2. `api/internal/profile/handler` & `repository` に一覧取得を追加（or article 側に追加）→ swag アノテーション
3. `cd api && make swagger` → `cd app && npm run generate:api`
4. `ProfileView.vue` をヘッダー + タブ構成に刷新
5. （論点3がAなら）アバターURL対応
