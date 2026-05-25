# 第1章 CASE式のススメ

> 学習日: 2026-05-25  |  関連まとめ: [summary.md#1章](../tatsujin-sql-shinansho-summary.md)

## 章の要旨（3-5行）

CASE は「文」ではなく **値を返す「式」** である。だから `SELECT` / `WHERE` / `GROUP BY` / `ORDER BY` / `UPDATE SET` / `CHECK制約` など、式が書ける場所ならどこにでも置ける。
この性質を使うと、(a) 集約軸の変換（`SELECT` と `GROUP BY` に同じ CASE）、(b) 行→列の水平展開（`SUM(CASE WHEN ... THEN 1 ELSE 0 END)` による条件付き集計＝クロス集計）、(c) 複数 UPDATE の1本化（順序依存バグの回避）、(d) CHECK 制約での「A ならば B」表現が、一時テーブルや手続きなしで宣言的に書ける。
落とし穴は、上から評価され最初にマッチで打ち切られる（順序）、戻り値の型を統一する、`ELSE` 省略は暗黙の `ELSE NULL`、NULL は `=` でなく `IS NULL` で書く（単純CASE式 `CASE col WHEN NULL` は永遠に真にならない）。

## このリポジトリでの該当箇所

| 種別 | ファイル:行 | 関連する論点 |
|---|---|---|
| migration | `api/migrations/20260518000000_change_replies_kind_to_varchar.up.sql:13` | **良い例**: 型変更 `ALTER COLUMN ... USING` 句で CASE による値変換（0→'comment', 1→'question', 2→'answer'）。「式が書ける場所ならどこでも」の好例 |
| Go クエリ | `api/internal/article/repository/get_list_articles.go:30` | いいね数を `LEFT JOIN (… COUNT(*) … GROUP BY) + COALESCE(…, 0)` で集計。CASE は未使用だが結果として妥当 |
| Go クエリ | `api/internal/reply/repository/get_list_replies.go` | 返信を kind 含め全件取得するのみ。集計・分類は行わない |
| TS 集計 | `app/src/features/replies/ReplySection.vue:30-38` | `counts` を JS で算出（`if (r.kind === 'question') c.question++`）。本章 (b) の `SUM(CASE WHEN kind='question' ...)` に相当する処理をフロントで実施 |

`api/` の **SELECT クエリ内に CASE 式は 0 件**。条件分岐が必要な集計（いいね数）は `LEFT JOIN + COUNT + COALESCE` で表現されており、現状の要件では CASE を使わずとも素直に書けている。

## 良い点（リポジトリで既に守れていること）

- migration の型変換で CASE を「式」として正しく活用している（手続き的な複数 UPDATE に頼っていない）。
- `get_list_articles.go` で `COALESCE(l.like_count, 0)` を使い、JOIN で生じる NULL を明示的に 0 へ畳んでいる。本章の「ELSE/NULL を放置しない」精神と整合。
- 集約 (`COUNT`, `array_agg`) を SQL 側で完結させ、フロントへ生データを撒いていない箇所が多い。

## 改善余地

### 案A: `ReplySection.vue` の kind 別カウントについて（結論: 据え置き推奨）
- 対象ファイル: `app/src/features/replies/ReplySection.vue:30-38`
- 観察: kind 別件数を JS の `for` ループで数えており、本章 (b) の `SUM(CASE WHEN kind=... THEN 1 ELSE 0 END)` 相当。
- 判断: この画面は **スレッド表示のために全返信をどのみち取得済み**なので、件数を別途 SQL 集計に切り出す利点は薄い（往復が増えるだけ）。「フロントで集計＝即悪」ではなく、**取得済みデータの再集計**ならフロントで問題ない。本章が問題視するのは「集計目的のためだけに生データを大量に持ってくる」パターン。よって現状維持を推奨。
- 学習ポイント: CASE/SUM による条件付き集計が活きるのは「集計結果だけが欲しい一覧系クエリ」。詳細画面のように元データが必要な場合は別。

### 案B: （改修候補なし）
- SELECT クエリに CASE を無理に差し込む箇所は現状ない。本章は **観点2（発展提案）で主役を立てる**のが学習効果が高い。

## このアプリの発展提案（新機能・集計アイデア）

> このアプリでは「質問」は記事にぶら下がる `kind='question'` の返信としてモデル化され、`answer` 返信のうち `is_best=true` がベストアンサー。
> articles に種別カラムは無い。この構造に対して CASE が主役になる集計を提案する。

### 発展提案1: 記事一覧の「質問ステータスバッジ」（未回答 / 回答あり / 解決済）
- 対応する PBI: No.3「多くの人が見てくれる場所で質問したい」/ No.4「最新記事を一覧で見たい」/ No.8「優れた回答に高評価」
- 何ができるようになるか: 記事一覧で各記事が抱える質問の状態（誰も答えていない / 回答はあるが未解決 / ベストアンサー確定）を一目で判別でき、未回答の質問に人を誘導できる。
- 本章テクニックがどう主役になるか: ステータスは「ベストアンサーがあるか」「回答返信があるか」「質問返信があるか」の **3段階の優先順位付き分類**。`CASE WHEN ... THEN ... WHEN ... THEN ... ELSE` の **上から評価・最初にマッチで打ち切り**という性質がそのまま仕様になる。素朴に3本のフラグを返してフロントで `if` 分岐すると、優先順位のロジックがクライアントに散る。CASE 1本ならステータス導出が SQL に一元化される。
- 必要なスキーマ・クエリの骨子（擬似SQL、新スキーマ不要）:
  ```sql
  -- get_list_articles.go の SELECT に列を1つ追加するイメージ
  SELECT a.id, a.title, /* ... 既存列 ... */
    CASE
      WHEN EXISTS (SELECT 1 FROM replies r
                   WHERE r.article_id = a.id AND r.is_best) THEN 'solved'
      WHEN EXISTS (SELECT 1 FROM replies r
                   WHERE r.article_id = a.id AND r.kind = 'answer') THEN 'answered'
      WHEN EXISTS (SELECT 1 FROM replies r
                   WHERE r.article_id = a.id AND r.kind = 'question') THEN 'unanswered'
      ELSE 'none'   -- 質問が付いていない通常記事
    END AS question_status
  FROM articles a;
  ```
- 規模感: 小〜中（クエリへ1列追加＋OpenAPI 型再生成＋一覧UIにチップ表示）。新テーブル不要。
- 既存実装との接続点: `api/internal/article/repository/get_list_articles.go` の SELECT、`article.Article` 構造体に `QuestionStatus` を追加、`app/src/features/articles/ArticlesView.vue` で Vuetify の `v-chip` 表示。

### 発展提案2: 投稿者ダッシュボードの「貢献サマリ」クロス集計
- 対応する PBI: No.6「自分の記事への評価を知りたい」/ No.9「投稿者・回答者のプロフィールを知りたい」
- 何ができるようになるか: あるユーザーが「コメント◯件 / 質問◯件 / 回答◯件 / うちベストアンサー◯件」を1行のクロス集計で把握できる（プロフィール画面の活動サマリ）。
- 本章テクニックがどう主役になるか: 1つの `replies` テーブルを kind ごとに **行→列へ水平展開**する典型 (b)。`SUM(CASE WHEN kind='answer' AND is_best THEN 1 ELSE 0 END)` のように **条件の合成**を列ごとに書けるのが CASE の真骨頂。kind ごとに別クエリを投げて JS で合算する素朴実装より、1スキャン1クエリで完結する。
- 必要なスキーマ・クエリの骨子（擬似SQL、新スキーマ不要）:
  ```sql
  SELECT
    SUM(CASE WHEN kind = 'comment'  THEN 1 ELSE 0 END) AS comment_count,
    SUM(CASE WHEN kind = 'question' THEN 1 ELSE 0 END) AS question_count,
    SUM(CASE WHEN kind = 'answer'   THEN 1 ELSE 0 END) AS answer_count,
    SUM(CASE WHEN kind = 'answer' AND is_best THEN 1 ELSE 0 END) AS best_answer_count
  FROM replies
  WHERE user_id = $1;
  ```
- 規模感: 中（新エンドポイント `GET /users/{id}/contribution` ＋ プロフィールUI）。新テーブル不要。
- 既存実装との接続点: `api/internal/profile/` に repository/handler を追加、`app/src/features/...` のプロフィール画面に統計カードを追加。

## 採用した対応

- 採用した改修案（観点1）: なし（SELECT に CASE を無理に差し込む箇所が無いため）
- 着手する発展提案（観点2）: **発展提案1（記事一覧の質問ステータスバッジ）** を縦切りで実装
- 実際に変更したファイル:
  - `api/internal/article/article.go` — `QuestionStatus` 型と定数、`Article.QuestionStatus` フィールド追加
  - `api/internal/article/repository/get_list_articles.go` — SELECT に優先順位付き CASE 式（solved > answered > unanswered > none）を追加し scan
  - `api/internal/article/handler/get_list_articles.go` — `ArticleListItemJSON.QuestionStatus`（swag `enums` 付き）とマッピング
  - `api/swagger/openapi.yml`, `api/swagger/swagger.json` — `make swagger` で再生成
  - `app/src/api/generated/apiSchema.ts` — `npm run generate:api` で再生成（`question_status: "none" | "unanswered" | "answered" | "solved"`）
  - `app/src/features/articles/ArticleCard.vue` — ステータス→チップの宣言的マップ + タイトル横バッジ表示（none は非表示）
- 補足: MSW モック (`app/src/mocks/db/articles.ts`) には `question_status` を付与していないため、モック起動時はバッジ非表示（`v-if` で安全に degrade）。実 API では機能する。

## 学びの言語化（自分用メモ）

- CASE は「式」。`USING` 句や `CHECK` 制約のような **SELECT 以外の場所**でも使えると意識すると、このリポジトリの migration（kind 型変換）のような実用例が見えてくる。
- 「フロントで集計しているから SQL に寄せるべき」は短絡。**取得済みデータの再集計はフロントで OK**、**集計だけが欲しいのに生データを撒くのが NG**、という線引きが本章の本質。
- このアプリで CASE が一番効くのは「優先順位付きステータス導出（質問の解決状態）」と「1テーブルの kind 別クロス集計（貢献サマリ）」。どちらも `replies` の `kind` / `is_best` という既存カラムだけで新スキーマなしに実現できる。
