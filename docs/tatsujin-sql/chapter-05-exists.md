# 第5章 EXISTS述語の使い方

> 学習日: 2026-05-25  |  関連まとめ: [summary.md#5章](../tatsujin-sql-shinansho-summary.md)

## 章の要旨（3-5行）

EXISTS は「サブクエリの結果が空でないか」を返す述語だが、背後には述語論理の**量化子（存在 ∃ / 全称 ∀）**がある。SQL に ∀ の構文は無いので、「すべての ◯◯ が P を満たす」を「**P を満たさない ◯◯ は存在しない**」と言い換え、`NOT EXISTS (... WHERE NOT P)` の**二重否定**で表す。EXISTS は **2値論理（TRUE/FALSE のみ、UNKNOWN を返さない）** で振る舞うため、`NOT IN` の NULL 消失罠を回避できる。サブクエリの先頭1行が見つかれば打ち切れるので早期終了で速いことが多い。SELECT リストは `SELECT *`/`1`/`NULL` どれでもよい。

## このリポジトリでの該当箇所

| 種別 | ファイル:行 | 関連する論点 |
|---|---|---|
| Go クエリ | `api/internal/article/repository/get_list_articles.go:26-31` | `EXISTS(...)` で「いいね済み判定」「回答ステータス(CASE+EXISTS)」。相関サブクエリだが NULL に強い |
| Go クエリ | `api/internal/article/repository/get_article.go:28` | 記事詳細の「いいね済み判定」も `EXISTS` |
| Go クエリ | `api/internal/reply/repository/set_best_answer.go:57` | `SELECT EXISTS(SELECT 1 ... is_best=TRUE)` で「既にベストアンサーがあるか」 |
| Go クエリ | `api/internal/like/repository/delete_like.go:29` | `SELECT EXISTS(SELECT 1 FROM articles WHERE id=$1)` で記事存在チェック → bool |
| Go クエリ | `api/internal/reply/repository/create_reply.go:101-105` | `ensureArticleExists`: `SELECT 1 FROM articles WHERE id=$1` + `ErrNoRows`。存在チェックだが EXISTS を使わない別スタイル |
| Go クエリ | `api/internal/like/repository/create_like.go:20-34` | 事前 EXISTS せず INSERT し、FK(23503)/UNIQUE(23505) 違反を捕捉。TOCTOU を避ける良い設計 |

**NOT EXISTS / 全称量化(∀) は現状コードに登場しない**。本章の真骨頂である「二重否定で ∀ を表す」テクニックは未活用で、発展提案の主戦場になる。

## 良い点（リポジトリで既に守れていること）

- **存在判定を `EXISTS` / `SELECT EXISTS(...)` で書けている**。記事一覧の「いいね済み」「回答ステータス」は相関サブクエリ + EXISTS で、`COUNT(*) > 0` より意図が明確かつ早期終了で効率的。5章慣習の `SELECT 1` も使えている。
- **EXISTS の 2値論理性を無意識に活かせている**。`liked_by_me` を `user_id NOT IN (...)` 等で書いていないため、4章で見た NULL 消失の罠に踏み込んでいない。
- **書き込みは「投機的存在チェック」より制約違反の捕捉を優先**（`create_like.go`）。「EXISTS で確認 → INSERT」の間に他トランザクションが割り込む TOCTOU を回避できる、実務的に正しい判断。

## 改善余地

### 案A: 存在チェックのスタイルを `SELECT EXISTS(...)` に統一する
- 対象ファイル: `api/internal/reply/repository/create_reply.go:101-105`（`ensureArticleExists`）
- 変更内容: 現状 `SELECT 1 FROM articles WHERE id=$1` を `Scan` して `sql.ErrNoRows` で分岐している。これを `delete_like.go` と同じ `SELECT EXISTS(SELECT 1 FROM articles WHERE id=$1)` → `bool` に揃える。
- 期待効果: 存在チェックの書き方がリポジトリ内で1スタイルに統一され、読み手の認知負荷が下がる。
- 学習ポイント: 「行の有無を知りたいだけなら値を取り出すより `EXISTS` で真偽を受ける」という5章の基本姿勢。ただし**現状でも正しく動いており、効果は可読性の統一に留まる**（章テクニックの主役性は弱め）。

> 注: EXISTS は既に適切に使えているため、観点1の改修インパクトは小さい。本章の学びは下の発展提案（NOT EXISTS / 全称量化）で主に得られる。

## このアプリの発展提案（新機能・集計アイデア）

### 発展提案1: 「まだいいねしていない人気記事」おすすめ（∃ の否定）
- 対応する PBI: No.1「記事の詳細を参照したい」/ No.4「最新の記事を一覧で見たい」（の発展＝発見導線）
- 何ができるようになるか: ログインユーザーに「自分がまだいいねしていない、いいねの多い記事」を数件おすすめする。既読・既いいねを除いた“次に読むべき記事”を出せる。
- 本章テクニックがどう主役になるか: 「自分のいいねが**存在しない**記事」を `NOT EXISTS` で表すのが核心。素朴に `a.id NOT IN (SELECT article_id FROM likes WHERE user_id=$1)` と書くと、サブクエリ側に NULL が紛れた瞬間に結果が全消失する（4章の罠）。`NOT EXISTS` なら 2値論理で安全、かつ早期終了で速い——5章の主張がそのまま効く。
- 必要なスキーマ・クエリの骨子:
  ```sql
  SELECT a.id, a.title,
         (SELECT COUNT(*) FROM likes WHERE article_id = a.id) AS like_count
  FROM articles a
  WHERE NOT EXISTS (
    SELECT 1 FROM likes l
    WHERE l.article_id = a.id AND l.user_id = $1   -- 自分のいいねが存在しない
  )
  ORDER BY like_count DESC, a.created_at DESC
  LIMIT 5;
  ```
- 規模感: 中（クエリ + エンドポイント + 一覧 UI）。まず backend だけの縦切り（1クエリ + 1エンドポイント）でレビュー可能。新スキーマ不要。
- 既存実装との接続点: `api/internal/article/repository/` に追加（`get_popular_articles.go` の隣）。`likes` の `EXISTS` 判定は `get_list_articles.go` に前例あり。

### 発展提案2: タグAND検索を「全称量化(∀)」の二重否定で書く別解
- 対応する PBI: No.7「分類別に記事・質問を閲覧したい」（既存のタグAND絞り込みの拡張・学習用別解）
- 何ができるようになるか: 機能としては現行のタグAND検索（指定タグを**すべて**持つ記事）と同じ。だが実装を `HAVING COUNT(DISTINCT)=N`（6章の関係除算）ではなく `NOT EXISTS` の二重否定で書き、「∀ = ¬∃¬」を体得する。
- 本章テクニックがどう主役になるか: 「指定タグを全部持つ記事」＝「**指定タグのうち、その記事に付いていないものが存在しない**」と言い換える。これが5章の「全称量化は二重否定で表す」の教科書通りの適用。現行の HAVING 版（`get_list_articles.go:53-61`）と並べると、同じ ∀ を集合カウントと述語論理の2通りで書けることが分かる。
- 必要なスキーマ・クエリの骨子:
  ```sql
  -- 要求タグ集合 wanted のうち、記事 a に付いていないものが存在しない記事
  SELECT a.id
  FROM articles a
  WHERE NOT EXISTS (
    SELECT 1 FROM unnest($2::text[]) AS wanted(name)
    WHERE NOT EXISTS (
      SELECT 1 FROM article_tags at JOIN tags t ON t.id = at.tag_id
      WHERE at.article_id = a.id AND LOWER(t.name) = wanted.name
    )
  );
  ```
- 規模感: 小〜中（既存クエリの別実装。挙動は同じなので比較検証が要る）。新スキーマ不要。
- 既存実装との接続点: `get_list_articles.go` のタグ絞り込み分岐。**現行 HAVING 版は動いている**ので、置き換えるより学習用に並記・ベンチして使い分けを判断するのが安全。

## 採用した対応

- 採用した改修案（観点1）: なし
- 着手する発展提案（観点2）: なし（提案1・2はメモに保持）
- なしを選んだ場合の理由: EXISTS は既にリポジトリ内で適切に使えており、観点1の改修は可読性の統一に留まる。本章の核心（NOT EXISTS / 全称量化）は発展提案として記録し、学習目的としては読了＋分析で十分と判断。
- 実際に変更したファイル: なし

## 学びの言語化（自分用メモ）

- このリポジトリは EXISTS をよく使えており、4章の NULL 罠を構造的に避けられている。一方で **NOT EXISTS（∃の否定）と全称量化（∀=¬∃¬）は未活用**で、ここに伸びしろがある。
- 「すべての◯◯が条件を満たす」「◯◯を全部含む」を見たら、まず **`NOT EXISTS (... WHERE NOT P)` の二重否定**に翻訳できないか考える。HAVING COUNT 版（6章）との2通りを持っておくと表現の引き出しが増える。
- 「除外」系（自分のいいね/既読を除く、回答が無い質問）は `NOT IN` ではなく `NOT EXISTS` を第一候補にする。NULL 安全・早期終了の両得。
</content>
