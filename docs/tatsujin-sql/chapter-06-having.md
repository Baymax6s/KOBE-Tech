# 第6章 HAVING句の力

> 学習日: 2026-05-25  |  関連まとめ: [summary.md#6章](../tatsujin-sql-shinansho-summary.md)

## 章の要旨（3-5行）

HAVING は GROUP BY の付属品ではなく「**集合の性質**」を問う SQL の中核機能。「**WHERE は要素（行）の性質、HAVING は集合（グループ）の性質**」と二分して捉える。実体1つに複数行が対応するなら、それは集合だから HAVING。代表技は **関係除算（指定アイテムを全部含むグループ＝バスケット解析）** を `HAVING COUNT(...) = (要求数)` で書くこと、`HAVING COUNT(*) = COUNT(col)` で「NULL を含まないグループ」を抽出すること、`COUNT(*) = SUM(CASE WHEN cond THEN 1 ELSE 0 END)` で「全要素が条件を満たす（全称量化）」を表すこと。学習のコツはベン図。

## このリポジトリでの該当箇所

| 種別 | ファイル:行 | 関連する論点 |
|---|---|---|
| Go クエリ | `api/internal/article/repository/get_list_articles.go:53-61` | タグAND絞り込みが `HAVING COUNT(DISTINCT LOWER(t.name)) = $3`。**関係除算そのもの**（指定タグを全部持つ記事）。6章の真骨頂を既に実装 |
| Go クエリ | `api/internal/article/repository/get_list_articles.go:34-36` | `SELECT article_id, COUNT(*) GROUP BY article_id` で記事ごとのいいね数を集約（集合→1行） |
| Go クエリ | `api/internal/article/repository/get_popular_articles.go:28-37` | `GROUP BY t.id,t.name,a.id,a.title` + `COUNT(l.id)` でタグ×記事のいいね数を集約し RANK |
| TS 集計 | `app/src/features/replies/ReplySection.vue:30-39` | ルート返信を kind 別に件数集計（all/comment/question）。GROUP BY 相当をフロントで実施 |

## 良い点（リポジトリで既に守れていること）

- **関係除算を HAVING COUNT で正しく書けている**（`get_list_articles.go:60`）。「指定タグを**すべて**持つ記事」を `WHERE ... = ANY(tags)` で要素を絞り → `GROUP BY article_id` で記事ごとの集合に → `HAVING COUNT(DISTINCT ...) = 要求タグ数` で集合の性質を判定、という6章の定石どおりの三段構え。`DISTINCT` で同名タグの重複カウントも防いでいる。
- **「要素の絞り込み（WHERE）」と「集合の判定（HAVING）」を正しく役割分担**している。タグ名一致は WHERE、件数一致は HAVING と、本章の二分法に沿っている。
- 集約とランキングを分離（CTE で集約 → 外側でフィルタ）しており、ウィンドウ関数を WHERE に直書きできない制約も理解した構造になっている（`get_popular_articles.go`）。

## 改善余地

### 案A: 返信の kind 別件数をフロント集計から SQL の GROUP BY に寄せる（低インパクト）
- 対象ファイル: `app/src/features/replies/ReplySection.vue:30-39`（+ 集計エンドポイント / DTO / OpenAPI）
- 変更内容: ルート返信を kind 別に数える `counts`（all/comment/question）をフロントの `for` ループで出している。これを `SELECT kind, COUNT(*) FROM replies WHERE article_id=$1 AND parent_id IS NULL GROUP BY kind` のような SQL 集約で返す。
- 期待効果: 「集合の件数」という集約をデータ層に置ける。
- 学習ポイント: 「実体（記事）1つに複数行（返信）が対応する＝集合」を GROUP BY で数える6章の基本。**ただし返信一覧は元々全件フェッチしており、表示中の配列を数えるだけのフロント集計は妥当**。SQL 化は API 往復が増えるだけで実益が薄く、むしろ現状維持が読みやすい。観点1としては「現状で問題なし」と判断するのが妥当。

> 注: タグAND検索で関係除算を既に実装できているため、観点1の改修インパクトは小さい。本章の学びは下の発展提案（集合の性質を問う集計）で主に得られる。

## このアプリの発展提案（新機能・集計アイデア）

### 発展提案1: 「盛り上がっている質問」（回答が N 件以上集まった質問）
- 対応する PBI: No.3「多くの人から答えが欲しい」/ No.2「回答したい」
- 何ができるようになるか: 「回答が一定数以上付いた質問」だけを抽出して“活発な議論”一覧を作る。逆に「回答が少ない質問」を運営が拾う用途にも使える。
- 本章テクニックがどう主役になるか: 1つの質問に複数の回答行がぶら下がる＝**質問は回答の集合**。その集合の大きさ（性質）を問うのが HAVING。`WHERE`（行のフィルタ）では「回答が N 件以上の質問」は表現できず、`GROUP BY 質問 HAVING COUNT(回答) >= N` が必要——「WHERE は要素・HAVING は集合」の二分法がそのまま効く。
- 必要なスキーマ・クエリの骨子:
  ```sql
  SELECT q.id, q.article_id, q.content, COUNT(a.id) AS answer_count
  FROM replies q
  LEFT JOIN replies a ON a.parent_id = q.id AND a.kind = 'answer'
  WHERE q.kind = 'question'
  GROUP BY q.id, q.article_id, q.content
  HAVING COUNT(a.id) >= $1          -- 集合（回答群）の大きさで絞る
  ORDER BY answer_count DESC;
  ```
- 規模感: 中（クエリ + エンドポイント + 一覧 UI）。まず backend の縦切り（1クエリ + 1エンドポイント）でレビュー可能。新スキーマ不要。
- 既存実装との接続点: `api/internal/reply/repository/`（自己結合は3章で導入済み）。記事一覧の回答ステータスチップとも親和。

### 発展提案2: 「全タグ制覇」バッジ（関係除算のユーザー応用＝全称量化）
- 対応する PBI: No.9「投稿者のプロフィールを知りたい」/ No.7「分類別に閲覧したい」の裏付け
- 何ができるようになるか: 「存在するすべてのタグについて、そのタグの記事を1本以上書いたことがあるユーザー」を“皆勤”として表彰する、といったゲーミフィケーション。一般化して「指定したタグ集合を全部カバーした投稿者」も出せる。
- 本章テクニックがどう主役になるか: これは `get_list_articles` のタグAND検索（関係除算）を**ユーザー軸に応用**したもの。「全タグを使った投稿者」＝「投稿に現れたタグの種類数 = タグ総数」で、`HAVING COUNT(DISTINCT tag_id) = (SELECT COUNT(*) FROM tags)` と書く。集合（投稿者が使ったタグ集合）が母集合と一致するかを HAVING で問う、関係除算＝全称量化の典型。
- 必要なスキーマ・クエリの骨子:
  ```sql
  SELECT a.user_id, COUNT(DISTINCT at.tag_id) AS used_tags
  FROM articles a
  JOIN article_tags at ON at.article_id = a.id
  GROUP BY a.user_id
  HAVING COUNT(DISTINCT at.tag_id) = (SELECT COUNT(*) FROM tags);
  ```
- 規模感: 小〜中（集計クエリ + エンドポイント。バッジ UI を足すなら中）。新スキーマ不要。
- 既存実装との接続点: `api/internal/profile/repository/` に集計を追加し、プロフィール画面のバッジに使う。タグAND検索（6章関係除算）の理解をそのまま転用できる。

## 採用した対応

- 採用した改修案（観点1）: 案A（返信 kind 別件数をフロント集計から SQL の GROUP BY に移管）
- 着手する発展提案（観点2）: なし（提案1・2はメモに保持）
- なしを選んだ場合の理由: 発展提案は次の候補として保持。今回は GROUP BY による集合の件数集計を実コードで体験する案Aを実施。
- 実装メモ: 別エンドポイントだと2回フェッチになるため、既存の返信一覧レスポンスに `counts` を追加する形にした。集計は `SELECT kind, COUNT(*) ... WHERE parent_id IS NULL GROUP BY kind`（本章の GROUP BY）。
- 実際に変更したファイル:
  - `api/internal/reply/reply.go`（`RootKindCounts` 型を追加）
  - `api/internal/reply/repository/count_root_replies.go`（新規。GROUP BY kind 集計）
  - `api/internal/reply/handler/get_list_replies.go`（レスポンスに `counts` を追加）
  - `api/swagger/openapi.yml` / `api/swagger/swagger.json`（`make swagger` で再生成）
  - `app/src/api/generated/apiSchema.ts`（`npm run generate:api` で再生成）
  - `app/src/features/replies/ReplySection.vue`（フロント集計 `counts` computed を撤去し、サーバ集計値を使用）
  - 検証: `gofmt` 差分なし / `go vet`・`go build`（reply）成功 / `npm run type-check` 通過 / 対象ファイルの eslint クリーン

## 学びの言語化（自分用メモ）

- このリポジトリは既に「関係除算（指定タグを全部持つ記事）」を `HAVING COUNT(DISTINCT)=N` で実装できている。6章の最重要パターンの実例として手元にあるのは強い。
- 「○件以上」「全部含む」「全員が満たす」を見たら、まず **GROUP BY でグループ（集合）を作り、その性質を HAVING で問う**形に翻訳する。WHERE で書こうとして詰まったら「それは行の性質か集合の性質か」を自問する。
- 関係除算（HAVING COUNT 版）と二重否定（NOT EXISTS 版・5章）は同じ ∀ の2通りの表現。両方を引き出しに持つと、可読性・性能で使い分けられる。
</content>
