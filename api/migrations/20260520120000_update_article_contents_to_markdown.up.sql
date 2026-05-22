-- 既存環境向けに、初期 seed の content を Markdown 化した内容に更新する。
-- seed (1777018911_seed_articles.up.sql) は履歴互換のため plain text のまま残し、
-- 本マイグレーションを適用した環境のみ Markdown 表示になる。
-- 同タイトルのユーザー投稿を巻き込まないよう、seed 由来の (title, user_id) 組で限定する。
UPDATE articles SET content = $md$# 神戸大学でのハッカソン体験記

先日、神戸大学で開催されたハッカソンに参加してきました。チームで 48 時間かけて Web アプリを作る、というシンプルなルールでしたが、得るものは想像以上でした。

## 当日のスケジュール

| 時間 | 内容 |
| --- | --- |
| Day1 09:00 | 開会式・テーマ発表 |
| Day1 11:00 | チームビルディング |
| Day2 18:00 | 発表 |

## 使った技術スタック

- フロント: Vue 3 + Vuetify
- バックエンド: Go + Gin
- DB: PostgreSQL
- インフラ: Docker Compose

ローカル起動は以下のコマンドで済むようにしておきました。

```bash
docker compose up -d
cd app && npm run dev
```

## チームで決めた最初のルール

> 1 PR 1 機能。Draft で出して早めに見せる。

このルールだけ守ったおかげで、終盤のコンフリクト解消で消耗せずに済みました。

## 学び

- 完璧な設計より、動くものを早く出して合わせる方が結果的に早い
- 詰まったらすぐ口に出す（Slack より口頭が圧倒的に速い）

来年も参加する予定です！
$md$ WHERE title = '神戸大学でのハッカソン体験記'
  AND user_id = (SELECT id FROM users WHERE name = 'admin');

UPDATE articles SET content = $md$# Goで作るREST API入門

Gin を使ったシンプルな REST API の作り方をまとめます。

## 最小構成

`main.go` に以下を書くだけで動きます。

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/articles", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"articles": []string{"hello", "world"},
		})
	})

	r.Run(":8080")
}
```

## JSON のバインド

POST されたボディを構造体に流し込みたいときは `ShouldBindJSON` を使います。

```go
type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func createArticle(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 保存処理...
}
```

## ミドルウェア

ロギングや CORS は `Use` で挟みます。

```go
r.Use(gin.Logger())
r.Use(cors.Default())
```

## 動作確認

```bash
curl -X POST http://localhost:8080/articles \
  -H 'Content-Type: application/json' \
  -d '{"title":"hello","content":"world"}'
```

## まとめ

- ルーティング・ミドルウェア・JSON バインドの 3 つを覚えれば最低限の API は組める
- バリデーションは `binding` タグでほぼ済む
$md$ WHERE title = 'Goで作るREST API入門'
  AND user_id = (SELECT id FROM users WHERE name = 'user01');

UPDATE articles SET content = $md$# Vue 3 + Vuetifyで学ぶフロントエンド開発

Composition API と Vuetify の組み合わせで、状態管理とコンポーネント分割をどう書くかを整理します。

## コンポーネント例

`<script setup>` を使うと、ボイラープレートが激減します。

```vue
<script setup lang="ts">
import { ref, computed } from 'vue'

const count = ref(0)
const doubled = computed(() => count.value * 2)
</script>

<template>
  <v-card class="pa-4">
    <v-card-title>{{ count }} (x2 = {{ doubled }})</v-card-title>
    <v-btn color="primary" @click="count++">+1</v-btn>
  </v-card>
</template>
```

## props と emits

型を `defineProps` に渡せば、テンプレート側で補完が効きます。

```ts
const props = defineProps<{
  articleId: number
}>()

const emit = defineEmits<{
  (e: 'select', id: number): void
}>()
```

## Vuetify のユーティリティ

レイアウト・余白は基本 Vuetify のユーティリティクラスで済ませます。

- `d-flex` / `ga-4` / `justify-center`
- `pa-4` / `mb-6`
- `bg-grey-lighten-4`

## まとめ

> Composition API + `<script setup>` + TypeScript の組み合わせが現状ベスト

慣れるまで戸惑いますが、Options API より型と再利用がはるかに楽になります。
$md$ WHERE title = 'Vue 3 + Vuetifyで学ぶフロントエンド開発'
  AND user_id = (SELECT id FROM users WHERE name = 'user02');

UPDATE articles SET content = $md$# PostgreSQLのマイグレーション管理

`golang-migrate` を使った DB マイグレーションの運用方針をまとめます。

## ファイル命名

タイムスタンプ + 内容で命名します。

```text
20260520123000_add_email_to_users.up.sql
20260520123000_add_email_to_users.down.sql
```

## up / down のペア

up と down を必ずペアで作ります。down が書けない場合は、その変更を分割するべきサインです。

```sql
-- up
ALTER TABLE users
  ADD COLUMN email VARCHAR(255) NOT NULL DEFAULT '';

CREATE UNIQUE INDEX users_email_idx ON users (email);
```

```sql
-- down
DROP INDEX IF EXISTS users_email_idx;
ALTER TABLE users DROP COLUMN IF EXISTS email;
```

## 破壊的変更は二段階で

列削除や NOT NULL 化など、旧コードを壊す変更は **新コードのデプロイ後** に流します。

| ステップ | 内容 |
| --- | --- |
| 1 | 追加系の up を流す |
| 2 | 新コードをデプロイ |
| 3 | 破壊的な up を流す |

## CI への組み込み

PR ごとに up → down → up を走らせて、down が壊れていないか検証します。

```yaml
steps:
  - run: make migrate-up
  - run: make migrate-down
  - run: make migrate-up
```

## まとめ

- 必ず down を書く
- 破壊的変更はデプロイと分ける
- CI で up/down の整合を見張る
$md$ WHERE title = 'PostgreSQLのマイグレーション管理'
  AND user_id = (SELECT id FROM users WHERE name = 'user03');

UPDATE articles SET content = $md$# Dockerで開発環境を統一する

`docker compose` を使ってチームの開発環境を統一する Tips です。

## 最小構成の docker-compose.yml

```yaml
services:
  db:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: password
      POSTGRES_DB: app
    ports:
      - "5432:5432"
    volumes:
      - db-data:/var/lib/postgresql/data

  api:
    build: ./api
    ports:
      - "8080:8080"
    depends_on:
      - db
    environment:
      DATABASE_URL: postgres://postgres:password@db:5432/app?sslmode=disable

volumes:
  db-data:
```

## ボリュームでデータを永続化

DB をコンテナ削除でリセットしないよう、`volumes` で永続化します。

```bash
docker compose down       # データは残る
docker compose down -v    # データも消える
```

## M1 / M2 Mac の platform 警告

amd64 イメージを使う場合は `platform` を明示します。

```yaml
services:
  legacy:
    image: some/amd64-only-image
    platform: linux/amd64
```

ただし、可能ならマルチアーキ対応の公式イメージを優先しましょう。

## よく使うコマンド

```bash
# 起動 (バックグラウンド)
docker compose up -d

# ログ追跡
docker compose logs -f api

# 特定のサービスだけ再ビルド
docker compose build api && docker compose up -d api

# DB に入る
docker compose exec db psql -U postgres -d app
```

## まとめ

- volumes でデータ消失を防ぐ
- platform 指定よりマルチアーキイメージを優先
- compose のコマンドはチームに共有する
$md$ WHERE title = 'Dockerで開発環境を統一する'
  AND user_id = (SELECT id FROM users WHERE name = 'user01');
