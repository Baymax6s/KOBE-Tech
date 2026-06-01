-- デモ用の追加記事。既存の 6 記事は触らず、状態のバリエーションを増やす。
--
-- 記事ごとに見せたい状態:
--   1. Git記事            … タグ2件 / いいね多 / コメント＋ベストアンサー付きQ&A＋4階層ネスト
--   2. TypeScript記事     … タグ3件（多タグ）/ ベストアンサー付きQ&A
--   3. 蔵書管理アプリ     … 制作物(PBI11) / タグ4件 / いいね最多 / ベスト未選択のQ&A（回答はあるが未確定）
--   4. CSS記事            … タグ2件 / いいね少 / 未回答の質問（回答ゼロ）
--   5. 学習習慣記事       … タグ無し / コメント無し（空スレッド）
--   6. 最新記事          … 投稿直後でいいね・コメントともにゼロ（一覧の先頭に出る新着）
--
-- content は KOBE-Tech の markdown-it 設定 (html:false / linkify:true / breaks:false) で
-- 描画できる記法のみを使う。

-- ============ 追加タグ ============
-- tags は LOWER(name) ユニークなので、大小文字を畳んだ重複チェックで冪等にする。
INSERT INTO tags (name, created_at, updated_at)
SELECT v.name, v.created_at, v.updated_at
FROM (VALUES
    ('Git',         '2026-05-06 09:00:00+09'::timestamptz, '2026-05-06 09:00:00+09'::timestamptz),
    ('TypeScript',  '2026-05-10 09:00:00+09'::timestamptz, '2026-05-10 09:00:00+09'::timestamptz),
    ('CSS',         '2026-05-18 09:00:00+09'::timestamptz, '2026-05-18 09:00:00+09'::timestamptz),
    ('個人開発',     '2026-05-14 09:00:00+09'::timestamptz, '2026-05-14 09:00:00+09'::timestamptz),
    ('制作物',       '2026-05-14 09:00:00+09'::timestamptz, '2026-05-14 09:00:00+09'::timestamptz),
    -- 学内カリキュラム特化タグ（Qiita/Zen の汎用タグとの差別化。神戸電子の課題語彙）
    ('C言語課題36',   '2026-05-13 09:00:00+09'::timestamptz, '2026-05-13 09:00:00+09'::timestamptz),
    ('C#課題24',     '2026-05-21 09:00:00+09'::timestamptz, '2026-05-21 09:00:00+09'::timestamptz),
    ('Java課題11',   '2026-05-23 09:00:00+09'::timestamptz, '2026-05-23 09:00:00+09'::timestamptz),
    ('プロジェクト管理', '2026-05-14 09:00:00+09'::timestamptz, '2026-05-14 09:00:00+09'::timestamptz),
    ('デザイン思考',   '2026-05-18 09:00:00+09'::timestamptz, '2026-05-18 09:00:00+09'::timestamptz)
) AS v(name, created_at, updated_at)
WHERE NOT EXISTS (
    SELECT 1 FROM tags t WHERE LOWER(t.name) = LOWER(v.name)
);

-- ============ 追加記事 ============
INSERT INTO articles (title, content, user_id, created_at, updated_at)
SELECT v.title, v.content, u.id, v.created_at, v.updated_at
FROM (VALUES
    ('Gitのrebaseとmergeを使い分ける',
     $md$# Gitのrebaseとmergeを使い分ける

チーム開発を始めると必ずぶつかるのが `rebase` と `merge` の使い分けです。どちらも「別のブランチの変更を取り込む」操作ですが、履歴の残り方が違います。

## ひとことで言うと

| 操作 | 履歴 | 向いている場面 |
| --- | --- | --- |
| `merge` | 分岐をそのまま残す | 共有ブランチへの取り込み |
| `rebase` | 一直線に並べ直す | 自分の作業ブランチの整理 |

## merge

`main` の変更を作業ブランチに取り込むときの基本形です。

```bash
git switch feature/login
git merge main
```

マージコミットが増えますが、いつ何を取り込んだかが履歴に残ります。

## rebase

作業ブランチのコミットを `main` の先頭に付け替えて、履歴を一直線にします。

```bash
git switch feature/login
git rebase main
```

> 共有済みのブランチを rebase すると履歴が書き換わって事故になります。**push 済みのブランチは rebase しない** を合言葉にしましょう。

## 迷ったときの方針

- 自分だけの作業ブランチ → `rebase` できれいに整える
- みんなが見る `main` や `develop` → `merge` で安全に取り込む

最初は merge だけでも十分です。慣れてきたら rebase を少しずつ使ってみてください。
$md$,
     '鈴木健一',
     '2026-05-06 10:00:00+09'::timestamptz,
     '2026-05-06 10:00:00+09'::timestamptz),

    ('TypeScriptの型で安全なAPIクライアントを作る',
     $md$# TypeScriptの型で安全なAPIクライアントを作る

フロントから API を叩くとき、レスポンスの形を型で固めておくと、後からの変更にぐっと強くなります。

## まずレスポンスの型を決める

```ts
type Article = {
  id: number
  title: string
  likeCount: number
}
```

## fetch をラップする

`fetch` の戻りは `any` 寄りになりがちなので、ジェネリクスで受け取る型を明示します。

```ts
async function getJson<T>(path: string): Promise<T> {
  const res = await fetch(path)
  if (!res.ok) throw new Error(`request failed: ${res.status}`)
  return res.json() as Promise<T>
}

const articles = await getJson<Article[]>('/api/articles')
```

## なぜ型を付けるのか

- `articles[0].titel` のようなタイポをエディタが即座に教えてくれる
- API の項目が増減したとき、影響範囲がコンパイルエラーで一覧できる

最初は面倒に感じますが、人数が増えるほど「型が仕様書になる」効きめが大きくなります。詳しくは公式の https://www.typescriptlang.org/ も覗いてみてください。
$md$,
     '高橋美咲',
     '2026-05-10 11:00:00+09'::timestamptz,
     '2026-05-10 11:00:00+09'::timestamptz),

    ('個人開発した蔵書管理アプリを公開します',
     $md$# 個人開発した蔵書管理アプリを公開します

積ん読が増えすぎたので、手持ちの本を管理する Web アプリを個人開発しました。フィードバックをもらえると嬉しいです！

## 作ったもの

- バーコードを入力すると書影とタイトルを自動取得
- 「読みたい / 読書中 / 読了」の3ステータスで管理
- 月ごとの読了数をグラフで表示

## 技術構成

| 領域 | 使ったもの |
| --- | --- |
| フロント | Vue 3 + Vuetify |
| バック | Go + Gin |
| DB | PostgreSQL |

KOBE-Tech とほぼ同じ構成にして、授業で学んだことをそのまま個人開発に持ち込みました。

## 工夫したところ

```go
// 外部APIのレスポンスが遅いことがあるので、取得済みの書影はDBにキャッシュする
func (r *Repository) cacheCover(isbn string, url string) error {
    // ...
}
```

> 完璧を目指すより、まず自分が毎日使うものを作る。これが続けるコツでした。

## これから

- 読んだ本の感想メモ機能
- 友達におすすめを共有する機能

改善案やツッコミ、気軽にコメントください！
$md$,
     '渡辺翔太',
     '2026-05-15 14:00:00+09'::timestamptz,
     '2026-05-15 14:00:00+09'::timestamptz),

    ('CSS Grid と Flexbox の使い分けメモ',
     $md$# CSS Grid と Flexbox の使い分けメモ

レイアウトを組むたびに「Grid と Flexbox どっち？」と迷うので、自分用に整理しました。

## ざっくり方針

- **1方向に並べる**（横一列 / 縦一列）→ Flexbox
- **2方向で格子に並べる**（行と列）→ Grid

## Flexbox の例

```css
.toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
}
```

## Grid の例

```css
.cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}
```

両方を組み合わせて、外枠は Grid・中身は Flexbox にすることも多いです。まだ手探りなので、定番の使い分けがあれば教えてほしいです。
$md$,
     '伊藤さくら',
     '2026-05-19 16:00:00+09'::timestamptz,
     '2026-05-19 16:00:00+09'::timestamptz),

    ('プログラミング学習を続けるための小さな習慣',
     $md$# プログラミング学習を続けるための小さな習慣

学習を始めて半年、なんとか続けられているので、自分が効果を感じた習慣をまとめます。技術記事ではないですが、誰かの参考になれば。

## 続けるために決めたこと

- 1日15分でもいいから毎日コードに触る
- 詰まったら30分で区切って質問する
- 「動いた」をスクショして記録する

## やめたこと

- 完璧に理解してから次へ進もうとすること
- 教材を増やしすぎること

最初の頃は1つの教材を最後までやり切るほうが、結局は近道でした。

## いま頑張っていること

毎日少しずつですが、このサイトに学んだことをアウトプットしていこうと思っています。同じ初心者の人、一緒に続けましょう！
$md$,
     '中村大輔',
     '2026-05-24 20:00:00+09'::timestamptz,
     '2026-05-24 20:00:00+09'::timestamptz),

    ('Postmanで始めるAPIデバッグ入門',
     $md$# Postmanで始めるAPIデバッグ入門

フロントを書く前に、API が期待どおり動くかを Postman で確かめておくと手戻りが減ります。投稿したばかりの新着記事です。

## まず1本リクエストを投げる

1. メソッドに `GET` を選ぶ
2. URL に `http://localhost:8080/api/articles` を入力
3. Send を押す

これでレスポンスの JSON が見られます。

## 認証が必要なとき

ログイン API で受け取ったトークンを、`Authorization` ヘッダーに `Bearer <token>` の形でセットします。

```text
Authorization: Bearer eyJhbGciOi...
```

## 環境変数を使う

ホスト名やトークンを環境変数にしておくと、ローカルと本番の切り替えが一発です。慣れてきたらコレクションにまとめて、チームで共有すると便利ですよ。
$md$,
     '鈴木健一',
     '2026-05-30 09:30:00+09'::timestamptz,
     '2026-05-30 09:30:00+09'::timestamptz),

    ('C言語課題36：ポインタのイメージを掴む',
     $md$# C言語課題36：ポインタのイメージを掴む

課題36でつまずきやすい「ポインタ」を、最小コードで整理します。

## ポインタは「住所」

変数の値そのものではなく、値が置いてある場所（アドレス）を持つのがポインタです。

```c
int x = 10;
int *p = &x;   // p は x の住所を持つ

printf("%d\n", x);    // 10  値そのもの
printf("%d\n", *p);   // 10  住所をたどって値を取り出す（デリファレンス）
```

## 配列名はほぼポインタ

```c
int a[3] = {10, 20, 30};
int *p = a;             // a は &a[0] と同じ

printf("%d\n", p[1]);   // 20  → *(p + 1) と同じ
```

> `p[i]` は `*(p + i)` の短い書き方。これさえ押さえれば配列とポインタは地続きです。

分からなくなったら「ポインタ＝住所」「`*`＝住所をたどる」に戻ってきてください。
$md$,
     '鈴木健一',
     '2026-05-13 10:00:00+09'::timestamptz,
     '2026-05-13 10:00:00+09'::timestamptz),

    ('C#課題24：例外処理でよくある詰まり',
     $md$# C#課題24：例外処理でよくある詰まり

課題24の例外処理で「どこまで catch すべきか」を整理します。

## 握りつぶさない

その場で回復できる例外だけ catch し、それ以外は上位に投げます。

```csharp
try
{
    var text = File.ReadAllText(path);
    Process(text);
}
catch (FileNotFoundException ex)
{
    Console.WriteLine($"ファイルが見つかりません: {ex.FileName}");
}
```

## やりがちなアンチパターン

```csharp
catch (Exception)
{
    // 全部受けると、本当のバグまで隠れてしまう
}
```

> 「catch するのは、その場で回復できる例外だけ」を合言葉に。
$md$,
     '鈴木健一',
     '2026-05-21 10:00:00+09'::timestamptz,
     '2026-05-21 10:00:00+09'::timestamptz),

    ('Java課題11：継承とインターフェースの使い分け',
     $md$# Java課題11：継承とインターフェースの使い分け

課題11でよく迷う2つを、ひとことで整理します。

## ざっくり対比

| 仕組み | 役割 |
| --- | --- |
| 抽象クラス | 共通の実装も配りたい親 |
| インターフェース | 「できること」の約束 |

## インターフェースの例

```java
interface Drawable {
    void draw();   // 約束だけ（実装なし）
}

class Circle implements Drawable {
    public void draw() {
        System.out.println("○");
    }
}
```

> 実装を共有したい → 抽象クラス / 型の約束だけ → インターフェース。
$md$,
     '高橋美咲',
     '2026-05-23 11:00:00+09'::timestamptz,
     '2026-05-23 11:00:00+09'::timestamptz)
) AS v(title, content, user_name, created_at, updated_at)
JOIN users u ON u.name = v.user_name
WHERE NOT EXISTS (
    SELECT 1 FROM articles a WHERE a.title = v.title AND a.user_id = u.id
);

-- ============ 記事 × タグ ============
-- 記事は (title, user_name) で一意に識別する。
INSERT INTO article_tags (article_id, tag_id, created_at)
SELECT a.id, t.id, v.created_at
FROM (VALUES
    ('Gitのrebaseとmergeを使い分ける',              '鈴木健一', 'Git',         '2026-05-06 10:00:00+09'::timestamptz),
    ('Gitのrebaseとmergeを使い分ける',              '鈴木健一', '開発環境',     '2026-05-06 10:00:00+09'::timestamptz),
    ('TypeScriptの型で安全なAPIクライアントを作る',  '高橋美咲', 'TypeScript',  '2026-05-10 11:00:00+09'::timestamptz),
    ('TypeScriptの型で安全なAPIクライアントを作る',  '高橋美咲', 'フロントエンド', '2026-05-10 11:00:00+09'::timestamptz),
    ('TypeScriptの型で安全なAPIクライアントを作る',  '高橋美咲', 'REST API',    '2026-05-10 11:00:00+09'::timestamptz),
    ('個人開発した蔵書管理アプリを公開します',        '渡辺翔太', '個人開発',     '2026-05-15 14:00:00+09'::timestamptz),
    ('個人開発した蔵書管理アプリを公開します',        '渡辺翔太', '制作物',       '2026-05-15 14:00:00+09'::timestamptz),
    ('個人開発した蔵書管理アプリを公開します',        '渡辺翔太', 'Go',          '2026-05-15 14:00:00+09'::timestamptz),
    ('個人開発した蔵書管理アプリを公開します',        '渡辺翔太', 'Vue',         '2026-05-15 14:00:00+09'::timestamptz),
    ('個人開発した蔵書管理アプリを公開します',        '渡辺翔太', 'プロジェクト管理', '2026-05-15 14:00:00+09'::timestamptz),
    ('CSS Grid と Flexbox の使い分けメモ',          '伊藤さくら', 'CSS',         '2026-05-19 16:00:00+09'::timestamptz),
    ('CSS Grid と Flexbox の使い分けメモ',          '伊藤さくら', 'フロントエンド', '2026-05-19 16:00:00+09'::timestamptz),
    ('CSS Grid と Flexbox の使い分けメモ',          '伊藤さくら', 'デザイン思考',  '2026-05-19 16:00:00+09'::timestamptz),
    ('Postmanで始めるAPIデバッグ入門',              '鈴木健一', 'REST API',    '2026-05-30 09:30:00+09'::timestamptz),
    -- 学内課題タグ付きの記事（差別化の柱3: 学内に最適化された分類）
    ('C言語課題36：ポインタのイメージを掴む',         '鈴木健一', 'C言語課題36',  '2026-05-13 10:00:00+09'::timestamptz),
    ('C#課題24：例外処理でよくある詰まり',           '鈴木健一', 'C#課題24',    '2026-05-21 10:00:00+09'::timestamptz),
    ('Java課題11：継承とインターフェースの使い分け',  '高橋美咲', 'Java課題11',  '2026-05-23 11:00:00+09'::timestamptz)
    -- 「プログラミング学習を続けるための小さな習慣」はタグ無し記事として残す。
) AS v(article_title, user_name, tag_name, created_at)
JOIN users u    ON u.name  = v.user_name
JOIN articles a ON a.title = v.article_title AND a.user_id = u.id
JOIN tags t     ON LOWER(t.name) = LOWER(v.tag_name)
ON CONFLICT (article_id, tag_id) DO NOTHING;
