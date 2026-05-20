-- 新規記事作成画面の「Markdown 記法」ボタンから参照する、
-- アプリ内シードのチートシート記事を admin 投稿として追加する。
-- 内容は KOBE-Tech の markdown-it 設定 (html:false / linkify:true / breaks:false)
-- で実際にレンダリングできる記法のみに絞る。脚注・KaTeX・絵文字など
-- 標準では動かない記法は載せない（「あるのに動かない」混乱を避けるため）。
INSERT INTO articles (title, content, user_id, created_at, updated_at)
SELECT v.title, v.content, u.id, v.created_at, v.updated_at
FROM (VALUES
    ('Markdown 記法チートシート', $md$# Markdown 記法チートシート

KOBE-Tech の記事は Markdown で書きます。このページでは、記事作成画面で使える主な記法をまとめています。手元で写経するつもりで眺めてみてください。

## 見出し

行頭の `#` の数で見出しレベルを決めます。

```markdown
# 見出し1
## 見出し2
### 見出し3
#### 見出し4
```

記事の中では `##` ～ `####` を使うと読みやすく、`#` は記事タイトルと重複しがちなので避けるのがおすすめです。

## 段落と改行

KOBE-Tech は **空行で段落を分ける** スタイルです（GitHub と同じ）。
行末で改行しても段落は分かれません。文章はゆったり空行で区切ってください。

```markdown
これは 1 つの段落です。
この行も同じ段落になります。

空行を挟むと、新しい段落になります。
```

## 強調

| 記法 | 表示例 |
| --- | --- |
| `**強調**` | **強調** |
| `*斜体*` | *斜体* |
| `~~打ち消し~~` | ~~打ち消し~~ |

## リスト

`-` か `*` を行頭に置くと箇条書きになります。半角スペース 2 つでネストできます。

```markdown
- リスト1
- リスト2
  - リスト2-1
  - リスト2-2
- リスト3
```

番号付きリストは `1.` から始めます。番号がずれていても、表示時には自動で振り直されます。

```markdown
1. 番号付きリスト1
2. 番号付きリスト2
3. 番号付きリスト3
```

## リンク

`[表示テキスト](URL)` の形で書きます。

```markdown
[KOBE-Tech のリポジトリ](https://github.com/Baymax6s/KOBE-Tech)
```

URL をそのまま貼っても自動的にリンクになります。

```markdown
https://github.com/Baymax6s/KOBE-Tech
```

## インラインコード

文中でコードや識別子を目立たせたいときは、バッククォートで囲みます。

```markdown
変数 `count` は `number` 型です。
```

## コードブロック

バッククォート 3 つで囲むとコードブロックになります。開きのバッククォートの直後に **言語名** を書くとシンタックスハイライトが効きます。

````markdown
```ts
const greet = (name: string) => `Hello, ${name}`
console.log(greet('KOBE'))
```
````

```ts
const greet = (name: string) => `Hello, ${name}`
console.log(greet('KOBE'))
```

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, KOBE")
}
```

```sql
SELECT id, title
FROM articles
WHERE user_id = $1
ORDER BY created_at DESC;
```

利用できる言語の一覧は、記事作成画面の **「対応言語」** ボタンから確認できます。`ts` / `js` / `py` などの短い別名（エイリアス）にも対応しています。

## 引用

行頭に `>` を置くと引用になります。引用の中で `>>` を使うと入れ子の引用にできます。

```markdown
> 引用文です。
> 複数行に渡って書けます。
>
>> 入れ子の引用も可能です。
```

> 引用文です。
> 複数行に渡って書けます。
>
>> 入れ子の引用も可能です。

## テーブル

`|` と `---` で表を書けます。`:` の位置でカラムごとの寄せ方を指定できます。

```markdown
| 言語       | 用途         | 補足           |
| ---------- | :----------: | -------------: |
| TypeScript |   フロント   |       Vue 3    |
| Go         | バックエンド |        Gin     |
| SQL        |      DB      |   PostgreSQL   |
```

| 言語       | 用途         | 補足           |
| ---------- | :----------: | -------------: |
| TypeScript |   フロント   |       Vue 3    |
| Go         | バックエンド |        Gin     |
| SQL        |      DB      |   PostgreSQL   |

## 水平線

行に `---` だけを書くと水平線になります。セクションの区切りに便利です。

```markdown
---
```

---

## 画像について

KOBE-Tech には **画像をアプリ内にアップロード・保存する機能はありません**。Qiita のようにエディタへ画像をドロップして埋め込む、といった使い方はできません。

外部にアップロードした画像の URL を `![alt](url)` で参照すれば、本文中に画像を表示することはできます。ただしリンク切れや権利関係には十分注意してください。

```markdown
![ロゴ](https://example.com/logo.png)
```

## このアプリで使えない記法

混乱を避けるため、KOBE-Tech では **使えない** 記法をまとめておきます。

- 画像のアプリ内アップロード（外部 URL 参照のみ可）
- 生 HTML タグ（`<div>` 等）— 表示されません
- 脚注 (`[^1]`)
- 目次の自動生成
- 数式（KaTeX / MathJax）
- 図（Mermaid / PlantUML）
- 絵文字ショートコード（`:smile:` 等）

これらが必要になった場合は、別途 issue で相談してください。

---

最初から全部使いこなす必要はありません。**見出し・段落・コードブロック** の 3 つだけでも、ぐっと読みやすい記事になります。まずは書き始めてみましょう！
$md$,
     'admin',
     '2026-05-20 12:00:00+09'::timestamptz,
     '2026-05-20 12:00:00+09'::timestamptz)
) AS v(title, content, user_name, created_at, updated_at)
JOIN users u ON u.name = v.user_name
WHERE NOT EXISTS (
    SELECT 1 FROM articles a WHERE a.title = v.title AND a.user_id = u.id
);
