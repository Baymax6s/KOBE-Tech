# 第2章 必ずわかるウィンドウ関数

> 学習日: 2026-05-25  |  関連まとめ: [summary.md#2章](../tatsujin-sql-shinansho-summary.md)

## 章の要旨（3-5行）

集約関数は「複数行→1行」に潰すが、**ウィンドウ関数は行を潰さず、その行の周辺集合に対する集約・順位を列として並べる**（情報保全性）。これにより、相関サブクエリでしか書けなかった行間比較・順位付けが1パスで書ける。
構文は `関数 OVER (PARTITION BY ... ORDER BY ... ROWS/RANGE ...)` の3要素。PARTITION BY が集合をカット（行は潰れない）、ORDER BY がカット内の順序、フレーム句が前後 N 行/値の範囲。
ランキング系（`RANK` / `DENSE_RANK` / `ROW_NUMBER` / `NTILE`）は ORDER BY 必須・フレーム不要。集約系（`SUM`/`AVG`...）はフレームで累積和・移動平均になる。ROWS（物理行数）と RANGE（値の範囲）の違いに注意。

## このリポジトリでの該当箇所

| 種別 | ファイル:行 | 関連する論点 |
|---|---|---|
| Go クエリ | `api/internal/article/repository/get_list_articles.go:30` | いいね数を `LEFT JOIN (… COUNT(*) GROUP BY)` で集計。順位は付けていない。ウィンドウ関数で順位列を足す発展の足場 |
| TS 集計 | `app/src/features/replies/ReplySection.vue:58-95` | `sort()` は返信スレッドの時系列並べ替え。ランキングではない |

**ウィンドウ関数（`OVER(...)`）の使用は api 全体で 0 件**。理由: 現状の要件はランキング・累積・行間比較を必要としていない（一覧は単純な `ORDER BY created_at DESC`、いいね数は件数集約のみ）。本章は **観点2（ランキング機能の新設）で主役を立てる**のが筋。

## 良い点（リポジトリで既に守れていること）

- 集計（いいね数）を SQL 側で完結させ、フロントへ生データを撒いていない。ウィンドウ関数を導入する際も同じ場所（`get_list_articles.go`）に乗せられる素地がある。
- フロントの `sort()` は「取得済みデータの表示順」に限定されており、ランキング（=母集合全体に対する順位）と混同していない。

## 改善余地

### 案A: （改修候補なし）
- 既存クエリにウィンドウ関数を差し込んで読みやすくなる箇所は現状ない。`get_list_articles.go` のいいね数集約は `LEFT JOIN + COUNT + COALESCE` で十分宣言的。`COUNT(*) OVER (...)` に置き換えても行が潰れない利点を活かす場面が無く、むしろ複雑化する。
- よって本章は観点1としては「該当の改修なし」。学習効果は観点2の新規ランキング機能で得る。

## このアプリの発展提案（新機能・集計アイデア）

> ランキングはウィンドウ関数の独擅場。素朴に書くと「各行ごとに自分より上位の件数を数える相関サブクエリ」になり、行数×行数で破綻する。

### 発展提案1: タグ別「人気記事 TOP3」（カテゴリ内ランキング）
- 対応する PBI: No.7「分類別に記事を閲覧したい」/ No.4「最新記事を一覧で見たい」/ No.6「自分の記事の評価を知りたい」
- 何ができるようになるか: 各タグ（カテゴリ）の中でいいね数上位3件を並べた「カテゴリ別の注目記事」セクション。1つの記事が複数タグに属していても、タグごとに別々に順位が付く。
- 本章テクニックがどう主役になるか: **「グループごとに順位を振り、各グループ上位 N 件だけ取り出す」**は `RANK() OVER (PARTITION BY tag ORDER BY like_count DESC)` が本領。PARTITION BY でタグごとに窓を切り、ORDER BY で順位付けする。ウィンドウ関数を知らないと「タグごとに上位を出す」ために相関サブクエリ（各記事に対し同タグでより多くいいねされた件数を数える）を書く羽目になり、可読性も性能も崩れる。**行を潰さずに順位列を付けられる**情報保全性がそのまま効く。
- 必要なスキーマ・クエリの骨子（擬似SQL、新スキーマ不要）:
  ```sql
  -- 既存の articles / likes / tags / article_tags だけで完結
  WITH ranked AS (
    SELECT
      t.id   AS tag_id,
      t.name AS tag_name,
      a.id   AS article_id,
      a.title,
      COUNT(l.id) AS like_count,
      RANK() OVER (
        PARTITION BY t.id
        ORDER BY COUNT(l.id) DESC, a.created_at DESC
      ) AS rank_in_tag
    FROM tags t
    JOIN article_tags at ON at.tag_id = t.id
    JOIN articles a      ON a.id = at.article_id
    LEFT JOIN likes l    ON l.article_id = a.id
    GROUP BY t.id, t.name, a.id, a.title
  )
  SELECT * FROM ranked WHERE rank_in_tag <= 3
  ORDER BY tag_name, rank_in_tag;
  ```
  ポイント: 集約 `COUNT(l.id)` の結果をそのまま `ORDER BY` に渡せる（ウィンドウ関数は GROUP BY 後に評価される）。`WHERE rank_in_tag <= 3` はウィンドウ関数を直接 WHERE に書けないのでサブクエリ/CTE で1段挟む（落とし穴）。
- 規模感: 中（新エンドポイント `GET /tags/popular-articles` ＋ トップページの「カテゴリ別注目記事」UI）。新テーブル不要。
- 既存実装との接続点: `api/internal/article/repository/` に新メソッド、`handler/handler.go` にルート追加、`app/src/features/articles/` にランキング表示コンポーネント。

### 発展提案2: 投稿者プロフィールの「月別投稿数の累積推移」
- 対応する PBI: No.9「投稿者・回答者のプロフィールを知りたい」/ No.6「自分の記事の評価を知りたい」
- 何ができるようになるか: プロフィール画面で「2026-01: 2件（累計2）/ 2026-02: 3件（累計5）…」のような活動の伸びを折れ線で見せる。
- 本章テクニックがどう主役になるか: 月ごとの件数を出すのは GROUP BY だが、**「その月までの累計」**は `SUM(月次件数) OVER (ORDER BY month ROWS UNBOUNDED PRECEDING)` という累積和フレームの典型 (a)。月次集計を一度フロントに渡してから JS で累積を計算する素朴実装より、SQL 1本で「月・月次・累計」の3列が揃う方が宣言的。
- 必要なスキーマ・クエリの骨子（擬似SQL、新スキーマ不要）:
  ```sql
  SELECT
    month,
    monthly_count,
    SUM(monthly_count) OVER (ORDER BY month
                             ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS cumulative
  FROM (
    SELECT date_trunc('month', created_at) AS month, COUNT(*) AS monthly_count
    FROM articles
    WHERE user_id = $1
    GROUP BY date_trunc('month', created_at)
  ) m
  ORDER BY month;
  ```
- 規模感: 中（新エンドポイント `GET /users/{id}/activity` ＋ プロフィールのグラフUI）。新テーブル不要。
- 既存実装との接続点: `api/internal/profile/`、`app/src/features/...` のプロフィール画面。

## 採用した対応

- 採用した改修案（観点1）: なし（既存クエリにウィンドウ関数を足して良くなる箇所が無いため）
- 着手する発展提案（観点2）: **発展提案1（タグ別人気記事TOP3）** を縦切りで実装
- 実際に変更したファイル:
  - `api/internal/article/popular_articles.go` — `PopularArticle` / `TagRanking` ドメイン型（新規）
  - `api/internal/article/repository/get_popular_articles.go` — CTE + `RANK() OVER (PARTITION BY tag ORDER BY COUNT(likes) DESC)` のクエリ、tag 境界でのグルーピング（新規）
  - `api/internal/article/handler/get_popular_articles.go` — DTO（swag）+ ハンドラ。TOP3 固定（`popularArticlesPerTag = 3`）（新規）
  - `api/internal/article/handler/handler.go` — `GET /tags/popular-articles` ルート追加
  - `api/swagger/openapi.yml`, `api/swagger/swagger.json` — `make swagger` で再生成
  - `app/src/api/generated/apiSchema.ts` — `npm run generate:api` で再生成（`tagsPopularArticlesList` / `ServerTagRankingJSONResponse` 等）
  - `app/src/features/articles/PopularArticlesByTag.vue` — タグ別ランキングの横スクロールUI（新規、v-slide-group）
  - `app/src/features/articles/ArticlesView.vue` — 絞り込み無し時に上記を見出し直下へ差し込み
- 補足: MSW モックには `/api/tags/popular-articles` ハンドラ未追加のため、モック起動時はセクション非表示（catch → rankings 空 → `v-if` で degrade）。実 API では機能する。検証: `go build`/`go vet`/`gofmt` クリーン、フロント `type-check`/`lint` パス。

## 学びの言語化（自分用メモ）

- ウィンドウ関数の見分け方: 「**順位を付けたい / 累計を出したい / 前後の行と比べたい**」が出てきたら相関サブクエリではなく `OVER(...)` を疑う。
- 「グループごとに上位 N 件」は `RANK()/ROW_NUMBER() OVER (PARTITION BY ...) ` + CTE で `WHERE rank <= N`。**ウィンドウ関数は WHERE に直接書けない**ので1段挟む、が定石。
- このアプリで一番効くのは「カテゴリ（タグ）内ランキング」。`likes` と `article_tags` という既存テーブルだけで、新スキーマなしに PBI No.7「分類別閲覧」を一段リッチにできる。
- `RANK`（同順位で番号飛ぶ） / `DENSE_RANK`（飛ばさない） / `ROW_NUMBER`（必ず連番）の使い分けを、TOP N の要件（同点をどう扱うか）に応じて選ぶ。
