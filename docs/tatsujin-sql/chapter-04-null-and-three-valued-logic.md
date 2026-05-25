# 第4章 3値論理とNULL

> 学習日: 2026-05-25  |  関連まとめ: [summary.md#4章](../tatsujin-sql-shinansho-summary.md)

## 章の要旨（3-5行）

SQL は 真 / 偽 / **不明 (UNKNOWN)** の3値論理で、NULL は「値が無い」を表すマーカーであって値ではない。だから NULL を含む比較・演算はほぼ UNKNOWN になり、`WHERE`/`HAVING` は UNKNOWN 行を弾く。典型的な罠が **`NOT IN` にサブクエリの NULL が混ざると結果が全消失**、**`<>` で除外すると NULL 行が漏れる**、**`COUNT(*)` と `COUNT(col)` の差**。回避の定石は `NOT EXISTS`・`IS NULL`・`COALESCE`。設計指針は「NOT NULL を基本に据え、NULL を増やさない」「『行が無い』『値が NULL』『空文字』を混同しない」。

## このリポジトリでの該当箇所

| 種別 | ファイル:行 | 関連する論点 |
|---|---|---|
| Go クエリ | `api/internal/profile/repository/update_bio.go:13-31` | `UPDATE ... WHERE user_id=$2` のみ。プロフィール行が無いと 0 行 → `ErrUserNotFound`。「行が無い」と「値が空」を混同（**潜在バグ**） |
| migration | `api/migrations/20260521003039_create_user_profiles.up.sql:14-17` | `INSERT ... SELECT ... WHERE bio IS NOT NULL`。当時 bio は全員 NULL だったため**0 行 INSERT**＝全ユーザーがプロフィール行なし |
| migration | `api/migrations/20260512042829_add_bio_to_users.up.sql:1-2` | `ADD COLUMN bio TEXT`（NOT NULL でもデフォルトでもない）。意味のない NULL を生む起点 |
| Go/migration | `user_profiles.bio` / `object_key` | NULL 許容だが、アプリは NULL と '' を同一視（下記「良い点」の COALESCE）。NULL が独自の意味を持たない＝22章「NULL撲滅」対象 |
| Go クエリ | `api/internal/profile/repository/get_profile.go:20` | `COALESCE(user_profiles.bio, '')` を `sql.NullString Bio` に scan。COALESCE 済みなので `Bio.Valid` は常に true（型と実態の不一致・軽微） |

## 良い点（リポジトリで既に守れていること）

- **`NOT IN (サブクエリ)` を一切使っていない**。記事一覧の「いいね済み判定」「回答ステータス」はすべて `EXISTS` で書かれており（`get_list_articles.go:26-32`）、4章で最も警告される「NOT IN + NULL で結果全消失」の罠を構造的に回避できている。
- **LEFT JOIN で生じる NULL を COALESCE で確実に既定値へ変換**している:
  - `get_list_articles.go:23-25`: `COALESCE(l.like_count, 0)` / `COALESCE(tag_ids, ARRAY[]::integer[])`。空集合集約が NULL を返す罠（23章）を 0・空配列に正規化。
  - `get_article.go:25-27`、`get_profile.go:20` も同様。
- **本当に nullable な列は `sql.NullInt64 / NullString / NullTime` で受け、`Valid` で分岐**している（`reply.parent_id`、`profile.UserProfile.CreatedAt` など）。`get_profile` の created/updated は LEFT JOIN で行が無いと NULL になるのを NullTime + `Valid` 判定でハンドラ側の `*time.Time`（JSON で null）に正しく落としている（`handler/get_profile.go:69-79`）。

## 改善余地

### 案A: `update_bio` を UPSERT にして「行が無い」と「値が空」を分離する（おすすめ）
- 対象ファイル: `api/internal/profile/repository/update_bio.go`
- 背景（重要）: `bio` は NULL 許容で追加され一度も seed されなかったため、`user_profiles` への移行 INSERT（`WHERE bio IS NOT NULL`）は **0 行**だった。登録エンドポイントも無いので、**現状どのユーザーにもプロフィール行が存在しない**。その結果 `UPDATE user_profiles ... WHERE user_id=$2` は常に 0 行 → `ErrUserNotFound` を返し、**誰も bio を保存できない**。
- 変更内容: UPDATE を UPSERT に変える。既存の一意制約 `uq_user_profiles_user_id` をそのまま使える:
  ```sql
  INSERT INTO user_profiles (user_id, bio, is_uploaded, created_at, updated_at)
  VALUES ($2, $1, FALSE, NOW(), NOW())
  ON CONFLICT (user_id)
  DO UPDATE SET bio = EXCLUDED.bio, updated_at = NOW();
  ```
  これで行が無ければ作成、あれば更新になり、`RowsAffected` ベースの `ErrUserNotFound` 分岐は不要になる（user 自体の存在確認が要るなら FK で担保される／別途 users を確認）。
- 期待効果: bio 保存機能が実際に動くようになる。「プロフィール行が存在しない（=まだ書いていない）」と「bio が空文字」を別物として正しく扱える。
- 学習ポイント: 4章の推奨方針「『未知』と『適用不能』を区別する」「存在の階層を意識する」の実践。**行の有無と値の NULL/空を取り違えると、エラーは出ないのに機能が動かない**という典型例。

### 案B: `bio` を `NOT NULL DEFAULT ''` 化して NULL を撲滅する
- 対象ファイル: 新 migration（`api/migrations/YYYYMMDDhhmmss_bio_not_null.{up,down}.sql`） + 整理として `get_profile.go:20` の COALESCE 除去・`profile.go` の `Bio sql.NullString` → `string` 化
- 変更内容: `UPDATE user_profiles SET bio='' WHERE bio IS NULL;` でバックフィル → `ALTER COLUMN bio SET DEFAULT ''` → `ALTER COLUMN bio SET NOT NULL`。
- 期待効果: 「NULL bio」という意味のない第3状態が消え、bio は常に文字列。COALESCE や NullString ラッパが不要になり読みやすくなる。
- 学習ポイント: 4章/22章「NOT NULL を基本に据える」「アプリが NULL と '' を区別しないなら NULL は不要」。ただし migration + 型整理を伴うので案A より規模が大きく、案A（行が作られる前提）とセットで効く。

## このアプリの発展提案（新機能・集計アイデア）

### 発展提案1: 「未回答の質問」フィード
- 対応する PBI: No.2「回答者として質問に回答したい」/ No.3「質問者として多くの人に答えてほしい」
- 何ができるようになるか: 「まだ回答（answer）が1件も付いていない質問」だけを集めた一覧を出す。回答者は貢献できる質問を見つけやすく、質問者は放置されにくくなる。
- 本章テクニックがどう主役になるか: 「回答が付いていない質問」を素朴に `WHERE id NOT IN (SELECT parent_id FROM replies WHERE kind='answer')` と書くと、**`parent_id` に NULL（ルート投稿）が混ざった瞬間に `<> NULL` が UNKNOWN になり結果が全消失する**——4章最大の罠そのもの。これを `NOT EXISTS` で書くのが正解、という4章の核心が主役になる。
- 必要なスキーマ・クエリの骨子:
  ```sql
  -- kind='question' で、自分にぶら下がる answer が1件も無い質問
  SELECT q.id, q.article_id, q.content, q.created_at
  FROM replies q
  WHERE q.kind = 'question'
    AND NOT EXISTS (
      SELECT 1 FROM replies a
      WHERE a.parent_id = q.id AND a.kind = 'answer'
    )
  ORDER BY q.created_at DESC;
  ```
- 規模感: 中（クエリ + エンドポイント + 一覧 UI）。新スキーマ不要。
- 既存実装との接続点: `api/internal/reply/repository/` に未回答質問取得を追加し、記事一覧の「回答ステータス」チップ（`get_list_articles.go` の question_status）とも親和。

### 発展提案2: 回答者の「ベストアンサー率」集計（0件・NULL の安全な扱い）
- 対応する PBI: No.8「優れた回答に高い評価をしたい」/ No.9「回答者のプロフィールを知りたい」
- 何ができるようになるか: ユーザープロフィールに「回答数 / うちベストアンサー数 / ベストアンサー率」を表示し、回答の信頼性を可視化する。
- 本章テクニックがどう主役になるか: 率＝`best数 / 回答数` は**回答 0 件のユーザーで 0 除算**になり、また `COUNT(col)` と `COUNT(*)` の差・空集合集約が NULL を返す挙動（4章/23章）を踏まえないと壊れる。`NULLIF(分母,0)` で 0 除算を NULL に逃がし、`COALESCE(..., 0)` で表示用に 0 へ戻す、という4章の処方箋が主役になる。
- 必要なスキーマ・クエリの骨子:
  ```sql
  SELECT u.id, u.name,
    COUNT(*) FILTER (WHERE r.kind = 'answer')                AS answer_count,
    COUNT(*) FILTER (WHERE r.kind = 'answer' AND r.is_best)  AS best_count,
    COALESCE(
      COUNT(*) FILTER (WHERE r.kind = 'answer' AND r.is_best)::numeric
      / NULLIF(COUNT(*) FILTER (WHERE r.kind = 'answer'), 0), 0
    ) AS best_rate
  FROM users u
  LEFT JOIN replies r ON r.user_id = u.id
  GROUP BY u.id, u.name;
  ```
- 規模感: 中（集計クエリ + エンドポイント + プロフィール UI）。新スキーマ不要。
- 既存実装との接続点: `api/internal/profile/repository/` に集計を追加し、プロフィール画面に指標を表示。

## 採用した対応

- 採用した改修案（観点1）: 案A（`update_bio` を UPSERT 化）
- 着手する発展提案（観点2）: なし（提案1・2はメモに保持）
- なしを選んだ場合の理由: まずは実害のある潜在バグ（誰も bio を保存できない）を低リスクで解消する案Aを優先。NULL撲滅の案Bと発展提案は次の候補としてメモに残す。
- 実際に変更したファイル:
  - `api/internal/profile/repository/update_bio.go`（UPDATE → `INSERT ... ON CONFLICT (user_id) DO UPDATE` の UPSERT 化。FK 違反 23503 のみ `ErrUserNotFound` にマップ）
  - 検証: `gofmt -l` 差分なし / `go vet ./internal/profile/...` クリーン / `go build ./internal/profile/...` 成功

## 学びの言語化（自分用メモ）

- このリポジトリは `NOT IN` を避け `EXISTS`、LEFT JOIN の NULL を `COALESCE` で正規化、nullable は `sql.Null*` で受ける——4章の処方箋をかなり守れている。NULL に強いコードの実例として参照できる。
- 一方 bio まわりは「**列を nullable で足す → seed しない → NULL のまま移行ロジックが空振り**」という、NULL が静かに事故を起こす教科書的な連鎖が起きていた。エラーは出ないのに機能が動かないのが NULL 起因バグの怖さ。
- 集計を増やすとき（率・平均・ランキング）は、まず「分母 0」「空集合」「対象 0 件のユーザー」を `NULLIF` / `COALESCE` でどう扱うかを先に決める癖をつける。
</content>
