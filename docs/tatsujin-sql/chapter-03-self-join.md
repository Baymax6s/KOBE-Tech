# 第3章 自己結合の使い方

> 学習日: 2026-05-25  |  関連まとめ: [summary.md#3章](../tatsujin-sql-shinansho-summary.md)

## 章の要旨（3-5行）

同じテーブルに別名を付けて結合する「自己結合」だけで、手続き型のループでやっていた処理を集合演算に置き換えられる。真価は **非等値結合（`<`, `>`, `<>`）** との組み合わせで、`<>` なら順列、`<` なら組合せ（`(a,b)` は出すが `(b,a)` は出さない）になる。`parent_id` のような自己参照列を持つテーブルでは「子の行と親の行」を1クエリで横並びにできる。ただしテーブルを2回走査するためコストは高く、ランキング等はウィンドウ関数（2章）の方が速いことも多い。

## このリポジトリでの該当箇所

| 種別 | ファイル:行 | 関連する論点 |
|---|---|---|
| migration | `api/migrations/20260507025654_create_replies.up.sql:11` | `parent_id INT` が `replies(id)` を自己参照。自己結合の前提となる構造そのもの |
| Go クエリ | `api/internal/reply/repository/set_best_answer.go:24-58` | 返信→親返信を **2回の QueryRow** に分けて取得。単一レベルの自己結合1本にまとめられる |
| Go クエリ | `api/internal/reply/repository/unset_best_answer.go:24-66` | 同上（best answer 解除側も同じ2クエリ構造） |
| TS 集計 | `app/src/features/replies/ReplySection.vue:150-163` | `questionAuthorByReplyId`: 各返信を親に突き合わせ「親が質問なら質問者ID」を引く=**自己結合をフロントで再実装**している |
| TS 集計 | `app/src/features/replies/ReplySection.vue:51-117` | スレッドツリー組み立て（childrenByParent / 祖先経路）。こちらは**多段=再帰**なので3章より21章（再帰CTE）の領域 |

## 良い点（リポジトリで既に守れていること）

- `replies` の自己参照 FK に `ON DELETE CASCADE` が付いており、親返信を消すと子も整合的に消える。スレッド構造が壊れない設計になっている。
- `idx_replies_parent_id` が貼られている。自己結合や子の引き当てで `parent_id` を辿る際にインデックスが効く土台がある。
- ツリーの「多段ネスト」をフロントで組み立てているのは、`ListByArticleID` で記事分の返信を1回で取り切ってから組むため N+1 を避けられており、判断として妥当（多段の自己結合/再帰をDBで毎回回すよりシンプル）。

## 改善余地

### 案A: best answer 操作の「返信→親」取得を自己結合1本にまとめる（おすすめ）
- 対象ファイル: `api/internal/reply/repository/set_best_answer.go` / `unset_best_answer.go`
- 変更内容: いまは「対象返信を取得」→「その `parent_id` で親返信を再取得」と **2回 QueryRowContext** している。これを `replies` を2回別名で結合する1クエリに統合する:
  ```sql
  SELECT
    a.kind         AS answer_kind,
    a.parent_id    AS parent_id,
    a.is_best      AS answer_is_best,
    q.kind         AS parent_kind,
    q.user_id      AS question_user_id
  FROM replies a
  LEFT JOIN replies q ON q.id = a.parent_id
  WHERE a.id = $1
  ```
  取得後に Go 側で「answer か」「親が question か」「質問者==操作者か」を判定する流れは現状と同じ。`LEFT JOIN` にするのは `parent_id` が NULL（ルート投稿）でも1行返し、Go 側の `ErrNotAnswer` 判定に乗せるため。
- 期待効果: DB 往復が2→1に減る。トランザクション内のクエリ本数が減り、ロジックも「1行を読んで分岐」に整理される。
- 学習ポイント: **同一テーブルを `a`（子=回答）/`q`（親=質問）の2別名で結合する**のが自己結合の最も基本の形。`parent_id` 自己参照テーブルに対する単一レベル結合の典型例で、API も UI も変えずに章のテクニックを体験できる。

### 案B: `questionAuthorByReplyId` の自己結合をバックエンドに寄せる
- 対象ファイル: `api/internal/reply/repository/get_list_replies.go`（+ DTO / OpenAPI）, `app/src/features/replies/ReplySection.vue:150-163`
- 変更内容: フロントが「各返信→親返信を引いて、親が質問なら質問者IDを記録」している処理は、`ListByArticleID` の SQL に自己結合を足せば DB 側で計算できる:
  ```sql
  SELECT r.id, ..., 
         CASE WHEN p.kind = 'question' THEN p.user_id END AS question_author_id
  FROM replies r
  JOIN users u ON u.id = r.user_id
  LEFT JOIN replies p ON p.id = r.parent_id
  WHERE r.article_id = $1
  ```
  返却 DTO に `question_author_id` を追加し、フロントの computed を削除する。
- 期待効果: 「誰がベストアンサーを付けられるか」の判定根拠がサーバ側に一元化され、フロントの集計コードが1つ減る。
- 学習ポイント: 自己結合 + CASE 式（1章）の合わせ技。ただし **OpenAPI 再生成を伴う**ぶん案Aより規模が大きく、フロントの既存挙動を壊さない検証も要る。

## このアプリの発展提案（新機能・集計アイデア）

### 発展提案1: 共通タグでつなぐ「関連記事」
- 対応する PBI: No.7「分類別に記事・質問を閲覧したい」/ No.1・No.4（記事を見つけたい・参照したい）
- 何ができるようになるか: 記事詳細ページの下部に「この記事と関連する記事」を、**共有しているタグ数が多い順**に数件表示する。回遊性が上がり、興味のある記事に辿り着きやすくなる。
- 本章テクニックがどう主役になるか: `article_tags` を**自己結合**し、`at1.tag_id = at2.tag_id AND at1.article_id <> at2.article_id` で「同じタグを共有する記事ペア」を作るのが核心。素朴にアプリ側で全記事のタグ配列を突き合わせると O(記事数²) のループになるが、自己結合 + GROUP BY なら宣言的に1クエリで書ける。
- 必要なスキーマ・クエリの骨子:
  ```sql
  -- 記事 :id と共有タグ数が多い順に関連記事を出す
  SELECT at2.article_id, COUNT(*) AS shared_tags
  FROM article_tags at1
  JOIN article_tags at2
    ON at1.tag_id = at2.tag_id
   AND at1.article_id <> at2.article_id   -- 自分自身は除く（順列の発想）
  WHERE at1.article_id = $1
  GROUP BY at2.article_id
  ORDER BY shared_tags DESC, at2.article_id
  LIMIT 5;
  ```
- 規模感: 中（クエリ追加 + エンドポイント + 記事詳細にUI）。新スキーマ不要。
- 既存実装との接続点: `api/internal/article/repository/` に `get_related_articles.go` を追加し、記事詳細ページ（記事取得まわりのコンポーネント）に関連記事カードを足す。

### 発展提案2: よく一緒に付けられるタグ（タグ共起サジェスト）
- 対応する PBI: No.7（分類で探す）の投稿側支援。投稿フォームでタグを選ぶときの体験向上。
- 何ができるようになるか: タグ A を選ぶと「A と一緒によく使われるタグ」を候補表示する。タグの表記ゆれ・付け忘れを減らせる。
- 本章テクニックがどう主役になるか: 同一記事内のタグ2つを**組合せ**として数えるのがそのまま自己結合。`t1.tag_id < t2.tag_id` で `{A,B}` を1回だけ数える（`(B,A)` を重複カウントしない）—3章の「組合せ＝`<`」の教科書通りの使い分けが主役になる。
- 必要なスキーマ・クエリの骨子:
  ```sql
  -- 同じ記事に同居したタグのペアを、共起回数が多い順に
  SELECT a.tag_id AS tag_a, b.tag_id AS tag_b, COUNT(*) AS co_count
  FROM article_tags a
  JOIN article_tags b
    ON a.article_id = b.article_id
   AND a.tag_id < b.tag_id              -- 組合せ（順序を排除）
  GROUP BY a.tag_id, b.tag_id
  ORDER BY co_count DESC;
  ```
- 規模感: 小〜中（共起クエリの追加が主。サジェストUIを足すなら中）。新スキーマ不要。
- 既存実装との接続点: `api/internal/article/repository/get_list_tags.go` の隣に共起集計を追加し、タグ入力コンポーネントの候補に流す。

## 採用した対応

- 採用した改修案（観点1）: 案A（best answer 操作の「返信→親」取得を自己結合1本に統合）
- 着手する発展提案（観点2）: なし（提案1・2はメモに保持）
- なしを選んだ場合の理由: 発展提案はいずれも価値があるが、まずは API/UI を変えず低リスクで自己結合の基本形を体験する案Aを優先。発展提案は次スプリント以降の候補としてメモに残す。
- 実際に変更したファイル:
  - `api/internal/reply/repository/set_best_answer.go`（2クエリ→自己結合 LEFT JOIN 1本）
  - `api/internal/reply/repository/unset_best_answer.go`（同上。`is_best` も同一クエリで取得）
  - 検証: `gofmt -l` 差分なし / `go vet ./internal/reply/...` クリーン / `go build ./internal/reply/...` 成功

## 学びの言語化（自分用メモ）

- 「同じテーブルを別名で2回出す」と身構えるより、**`parent_id` 自己参照テーブルは "子の行と親の行を横に並べたい" 場面で自然に自己結合になる**、と覚えるのが実戦的。best answer 判定の2クエリ往復はその典型で、1本に畳める。
- 「ペアを作る」系の機能（関連記事・共起タグ・おすすめユーザー）はすべて自己結合 + 非等値結合で書ける。`<>` か `<` の使い分け（順列 vs 組合せ）を意識すると重複カウント事故を防げる。
- 多段ツリー（スレッドの全階層）は3章の単純自己結合では届かず、**21章の再帰CTE（`WITH RECURSIVE`）**の領域。いまフロントでツリーを組んでいるのは妥当な判断で、DB側でやるなら再帰CTEを使う、という線引きを押さえておく。
</content>
</invoke>
