# 達人に学ぶSQL徹底指南書 第2版 まとめ

- **著者**: ミック
- **出版**: 翔泳社（CodeZine BOOKS） / ISBN 9784798157825
- **位置づけ**: 2008年初版の10年ぶり改訂。**ウィンドウ関数を全面導入**してコードを刷新。標準SQL準拠で Oracle / SQL Server / DB2 / PostgreSQL / MySQL などに対応。
- **対象読者**: 文法は分かる初級者が「中級者」になるための本。**SQLらしい思考（集合指向・宣言型）**を身につけることが目的。
- **構成**: 第1部 魔法のSQL（実践テクニック12章）/ 第2部 RDBの世界（理論11章）/ 第3部 付録。

> 本ファイルは Web 上の公開情報（出版社書誌、著者ミック氏のWeb連載、各種書評記事）から構成した二次的な学習メモ。SQLコードは典型パターンを示したもので、原書本文の完全な引用ではありません。

---

## 目次

### 第1部 魔法のSQL
1. CASE式のススメ
2. 必ずわかるウィンドウ関数
3. 自己結合の使い方
4. 3値論理とNULL
5. EXISTS述語の使い方
6. HAVING句の力
7. ウィンドウ関数で行間比較を行なう
8. 外部結合の使い方
9. SQLで集合演算
10. SQLで数列を扱う
11. SQLを速くするぞ
12. SQLプログラミング作法

### 第2部 リレーショナルデータベースの世界
13. RDB近現代史
14. なぜ"関係"モデルという名前なの?
15. 関係に始まり関係に終わる
16. アドレス、この巨大な怪物
17. 順序をめぐる冒険
18. GROUP BYとPARTITION BY
19. 手続き型から宣言型・集合指向へ頭を切り替える7箇条
20. 神のいない論理
21. SQLと再帰集合
22. NULL撲滅委員会
23. SQLにおける存在の階層

### 第3部 付録
- A 演習問題の解答
- B 参考文献

---

# 第1部 魔法のSQL

## 1章 CASE式のススメ

### この章の主張
**SQLにおける条件分岐の王様は IF文ではなく CASE式**。CASE は文（statement）ではなく "値を返す式（expression）" なので、`SELECT` / `WHERE` / `GROUP BY` / `ORDER BY` / `UPDATE SET` / `CHECK制約` など、**式が書ける場所ならどこにでも書ける**。これが他の条件分岐構文（DECODE, IF 関数）と比べた最大の強み。

### 2つの構文
- **単純CASE式**: `CASE col WHEN '1' THEN '男' WHEN '2' THEN '女' ELSE 'その他' END`
- **検索CASE式**: `CASE WHEN sex='1' THEN '男' WHEN sex='2' THEN '女' ELSE 'その他' END`

検索CASE式の方が表現力が高い（範囲条件、複合条件、サブクエリも書ける）。基本的にはこちらを使う。

### 典型テクニック

**(a) 既存コード値の集約軸を変換しながら GROUP BY**
```sql
SELECT
  CASE pref_name
    WHEN '徳島' THEN '四国' WHEN '香川' THEN '四国'
    WHEN '愛媛' THEN '四国' WHEN '高知' THEN '四国'
    WHEN '福岡' THEN '九州' WHEN '佐賀' THEN '九州'
    ELSE 'その他'
  END AS district,
  SUM(population) AS pop
FROM PopTbl
GROUP BY CASE pref_name WHEN '徳島' THEN '四国' ... END;
```
ポイント: **`SELECT` と `GROUP BY` に同じ CASE 式を書く**。一時テーブルを作らずに集約軸を作れる。

**(b) 行 → 列の水平展開（SUM/CASE による条件付き集計）**
```sql
SELECT pref_name,
  SUM(CASE WHEN sex = '1' THEN population ELSE 0 END) AS male,
  SUM(CASE WHEN sex = '2' THEN population ELSE 0 END) AS female
FROM PopTbl2
GROUP BY pref_name;
```
クロス集計表が1クエリで作れる。`SUM(CASE WHEN ... THEN 1 ELSE 0 END)` で**条件を満たす行数のカウント**にもなる。

**(c) 複数 UPDATE を1本化**
```sql
UPDATE Personnel SET salary =
  CASE WHEN salary >= 300000 THEN salary * 0.9
       WHEN salary >= 250000 AND salary < 280000 THEN salary * 1.2
       ELSE salary
  END;
```
**素朴に複数 UPDATE を書くと、1回目の更新で条件が変わって2回目が誤動作する**（順序依存バグ）。CASE 1本ならアトミック。

**(d) CHECK 制約での複合条件**
```sql
CONSTRAINT check_salary CHECK (
  CASE WHEN sex = '2'
       THEN CASE WHEN salary <= 200000 THEN 1 ELSE 0 END
       ELSE 1 END = 1
)
```
「女性社員の給与は20万円以下」のような**「A ならば B」**を CASE で表現する。論理学の含意 (A→B ≡ ¬A∨B) と等価。

**(e) テーブル間マッチング**
```sql
SELECT keyCol,
  CASE WHEN keyCol IN (SELECT keyCol FROM tbl_B) THEN 'Match'
       ELSE 'Unmatch' END AS label
FROM tbl_A;
```

### 落とし穴
- **CASEは上から評価され、最初にマッチした WHEN で打ち切られる**。順序を間違えると下の条件が永遠に評価されない。
- 全 WHEN/ELSE の**戻り値の型を統一する**（混在すると暗黙変換で事故る）。
- **ELSE を省略すると暗黙の `ELSE NULL`**。`SUM(CASE ...)` などで NULL が混ざると意図しない結果になりやすいので、原則明示する。
- **NULL の比較は `=` ではなく `IS NULL` で**書く。単純CASE式 `CASE col WHEN NULL THEN ...` は永遠に真にならない（`col = NULL` は UNKNOWN）。

---

## 2章 必ずわかるウィンドウ関数

### この章の主張
SQL の集約関数は「複数行 → 1行」に潰してしまうが、**ウィンドウ関数は「行を潰さずに、その行の周辺集合に対する集約結果を列として並べる」**。これにより、相関サブクエリでしか書けなかった行間比較が一発で書ける。本書第2版の **看板トピック**。

### ウィンドウの3要素
```sql
集約関数 / ランキング関数 OVER (
  PARTITION BY col1, col2   -- ① 集合をカット
  ORDER BY col3              -- ② カット内で順序づけ
  ROWS BETWEEN m PRECEDING AND n FOLLOWING  -- ③ フレーム（部分集合）
)
```
1. **PARTITION BY**: 全体をグループに分割（GROUP BY と違い、行は潰れない）。
2. **ORDER BY**: パーティション内での順序を決める（時系列・順位付けで必須）。
3. **フレーム句**: カレント行から見た **前後 N 行 / N 値の範囲** を切り出す。集約系（SUM/AVG/MAX...）で意味を持つ。

### 主な関数の分類
- **集約系**: `SUM`, `AVG`, `COUNT`, `MAX`, `MIN` を `OVER` 付きで使う。
- **ランキング系**: `RANK`, `DENSE_RANK`, `ROW_NUMBER`, `NTILE` — ORDER BY 必須、フレーム句は不要。
- **行間参照系**: `LAG`, `LEAD`, `FIRST_VALUE`, `LAST_VALUE` — 7章で詳しく扱う。

### 典型パターン

**(a) 累積和**
```sql
SELECT order_date, amount,
  SUM(amount) OVER (ORDER BY order_date
                    ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS cum_sum
FROM Sales;
```

**(b) 直近3行の移動平均**
```sql
SELECT order_date, amount,
  AVG(amount) OVER (ORDER BY order_date
                    ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) AS moving_avg
FROM Sales;
```

**(c) カテゴリ内ランキング（PARTITION BY と組み合わせ）**
```sql
SELECT product, category, price,
  RANK() OVER (PARTITION BY category ORDER BY price DESC) AS price_rank
FROM Products;
```

### ROWS と RANGE の違い（重要）
- **ROWS**: 物理的な行数で範囲を切る。「直前2行」など。
- **RANGE**: ORDER BY 列の**値**で範囲を切る。同値の行はまとめて扱われる。

例えば日付に重複がある時系列で `ROWS BETWEEN 2 PRECEDING ...` と `RANGE BETWEEN INTERVAL '2' DAY PRECEDING ...` は別物。

### コラム
- **「なぜ ON ではなく OVER なのか?」**: ウィンドウ関数は結合（JOIN）ではなく「窓を覗く」イメージ。元の行はそのままに、別軸の集計結果を覗き見るから OVER。

### 嬉しさ
- 相関サブクエリ撲滅 → 可読性 & パフォーマンス両改善。
- **情報保全性**: 元テーブルの行を1行も削らずに、列として情報を追加できる。

---

## 3章 自己結合の使い方

### この章の主張
**同じテーブルに別名を付けて結合する**だけで、ループや手続き型コードでやっていた処理を集合演算に置き換えられる。さらに **非等値結合（`<`, `>`, `<>`）** と組み合わせることで真価を発揮する。

### キーとなる発想: 「対」の数学
- **順序対** `<a, b>`: 順番が違えば別物。「a→b 移動」など方向のある情報に使う。
- **非順序対 / 組合せ** `{a, b}`: 順番無視。「商品 a と b のペア」など。

### 典型パターン

**(a) 重複順列（直積）**
```sql
SELECT P1.name AS n1, P2.name AS n2
FROM Products P1 CROSS JOIN Products P2;
-- (りんご,りんご),(りんご,みかん),(みかん,りんご),...
```

**(b) 順列（重複排除、順序あり）**
```sql
SELECT P1.name, P2.name
FROM Products P1 INNER JOIN Products P2
  ON P1.name <> P2.name;
```

**(c) 組合せ（重複も順序も排除）**
```sql
SELECT P1.name, P2.name
FROM Products P1 INNER JOIN Products P2
  ON P1.name < P2.name;
```
`<` を使うと **(a,b) のペアは出るが (b,a) は出ない**。これが組合せ。

**(d) 自己結合によるランキング（ウィンドウ関数なしでの実装）**
```sql
SELECT P1.name, P1.price,
  (SELECT COUNT(P2.price) FROM Products P2
    WHERE P2.price > P1.price) + 1 AS rank
FROM Products P1;
```
「自分より高い価格の商品が何個あるか」を数えれば、それが順位 -1。ウィンドウ関数登場前はこれが定石だった。

**(e) 重複行の削除（連番なしテーブル）**
```sql
DELETE FROM Products P1
WHERE rowid < (SELECT MAX(P2.rowid) FROM Products P2
                WHERE P1.name = P2.name AND P1.price = P2.price);
```

### コラム: SQLとフォン・ノイマン
手続き型言語の祖フォン・ノイマンと、集合指向のSQL。発想がほぼ正反対であることが、SQL を難しく感じる原因。

### 注意
**自己結合は強力だがコストが高い**。テーブル全体を2回スキャンするため、大きなテーブルでは性能を必ず計測する。多くの場面で**ウィンドウ関数の方が速い**ので、第2版以降ではウィンドウ関数版を優先するのが定石。

---

## 4章 3値論理とNULL

### この章の主張
SQL は **真 (TRUE) / 偽 (FALSE) / 不明 (UNKNOWN)** の3値論理。NULL は「値が無い」を表すマーカーであって**値ではない**。だから NULL を含む比較・演算はほぼ全て UNKNOWN になり、それが原因で「実行してもエラーにならないが結果が間違っている」クエリが量産される。

### 3値論理の真理表（要点）
- `1 = NULL`, `1 <> NULL`, `1 > NULL`, **`NULL = NULL`** … すべて UNKNOWN。
- `TRUE AND UNKNOWN` = UNKNOWN, `FALSE AND UNKNOWN` = FALSE。
- `TRUE OR UNKNOWN` = TRUE, `FALSE OR UNKNOWN` = UNKNOWN。
- **WHERE / HAVING / CHECK は「TRUE の行だけ」を通す**。UNKNOWN は弾かれる。

### よくある罠

**(a) `NOT IN` に NULL が混ざると結果が全部消える**
```sql
-- Class_A: 1, 2, 3, NULL
SELECT * FROM Students WHERE age NOT IN (SELECT age FROM Class_A);
-- → 結果が必ず空集合。`age <> NULL` が UNKNOWN になるため。
```
**回避**: `NOT EXISTS` を使う、もしくはサブクエリ側で `WHERE age IS NOT NULL` を付ける。

**(b) `<>` での除外が NULL を漏らす**
```sql
SELECT * FROM Tbl WHERE col <> 1;
-- → col が NULL の行は出てこない。
```
意図的に NULL も拾いたければ `WHERE col <> 1 OR col IS NULL`。

**(c) 集約関数は NULL を無視するが、COUNT(*) だけは違う**
- `COUNT(col)` … NULL を除外して数える。
- `COUNT(*)` … 行数そのもの（NULL も数える）。
- `AVG(col)` … 分母は NULL でない行数。**「分母に NULL を含めたい」場合は事前に COALESCE で 0 を入れるなど工夫**が必要。

**(d) `CASE col WHEN NULL THEN ...` は永遠に真にならない**
単純CASE式は内部で `=` 比較するため、NULL とは一致しない。検索CASE式で `WHEN col IS NULL THEN ...` と書く。

### 推奨方針
- **NOT NULL 制約を基本に据える**。NULL を許可するのは設計上の判断で例外的に。
- **デフォルト値で代用**できないか先に検討（`UNKNOWN` 区分コード, `9999-12-31` など）。
- 「未知 (unknown)」と「適用不能 (not applicable)」を区別したければ、別テーブルや明示フラグで設計する。

### コラム: 文字列とNULL
Oracle は**空文字列 `''` を NULL として扱う**という独自仕様。他DBと挙動が異なるので移植時に要注意。

---

## 5章 EXISTS述語の使い方

### この章の主張
EXISTS は単なる「サブクエリの結果が空でないか」を返す述語に見えるが、その背後には**述語論理の量化子（存在量化 ∃ と 全称量化 ∀）** が隠れている。SQL には ∀ に対応する構文がない（標準には `ALL` 述語はあるが弱い）ので、**`NOT EXISTS + NOT ...` の二重否定で全称量化を表現**する。

### 述語論理との対応
- **存在量化** ∃x P(x) ⇔ `EXISTS (SELECT ... WHERE P)`
- **全称量化** ∀x P(x) ⇔ ¬∃x ¬P(x) ⇔ `NOT EXISTS (SELECT ... WHERE NOT P)`

「すべての ◯◯ が条件 P を満たす」を「**条件 P を満たさない ◯◯ は存在しない**」と言い換えるのが鍵。

### 典型パターン

**(a) 全教科で50点以上の学生**
```sql
SELECT DISTINCT student_id FROM TestScores TS1
WHERE NOT EXISTS (
  SELECT * FROM TestScores TS2
  WHERE TS2.student_id = TS1.student_id
    AND TS2.score < 50
);
```
「50点未満の科目を1つも持たない学生」と読み替える。

**(b) 歯抜けのある連番を検出（欠番）**
```sql
SELECT seq FROM SeqTbl S1
WHERE NOT EXISTS (
  SELECT * FROM SeqTbl S2 WHERE S2.seq = S1.seq + 1
);
-- 自分の次の番号が存在しない行 = 抜けの直前
```

### EXISTS の特殊性
- **EXISTS は2値論理で振る舞う**（TRUE / FALSE のみ。UNKNOWN を返さない）ので、`NOT IN` のように NULL で結果が消失する罠を回避できる。
- **SELECT リストは何でもよい**: `SELECT * / SELECT 1 / SELECT NULL` どれでも同じ。本書では慣習として `SELECT *` 推奨。
- **EXISTS は高階の述語**: 述語を引数に取る、関数型言語の高階関数に近い性質を持つ。

### 性能面
`IN` よりも `EXISTS` の方が、サブクエリの最初の1行が見つかった時点で打ち切れるため**早期終了で速い**ことが多い。ただし、最近のオプティマイザは `IN` を内部的に EXISTS に書き換えるので、必ずしも書き分けで性能が変わるわけではない。

---

## 6章 HAVING句の力

### この章の主張
HAVING は GROUP BY の付属品ではなく、**「集合の性質」を問うための SQL の中核機能**。「WHERE は要素の性質、HAVING は集合の性質」と二分して捉えるのが本書の流儀。

### WHERE と HAVING の対比
| 観点 | WHERE | HAVING |
|------|-------|--------|
| 対象 | 個別の行（要素） | グループ（集合） |
| 使える式 | 列・スカラー値 | 集約関数・グループ列 |
| 思考モデル | 要素1個ずつフィルタ | 集合全体の性質を判定 |

「実体1つに複数行が対応している → それは集合だから HAVING」というのが見極めの方針。

### 典型パターン

**(a) 中央値**
HAVING で「自分以下の行数と自分以上の行数の差が小さい値」を捉える。
```sql
SELECT AVG(DISTINCT income) AS median FROM (
  SELECT T1.income FROM Graduates T1, Graduates T2
  GROUP BY T1.income
  HAVING SUM(CASE WHEN T2.income >= T1.income THEN 1 ELSE 0 END)
           >= COUNT(*) / 2
     AND SUM(CASE WHEN T2.income <= T1.income THEN 1 ELSE 0 END)
           >= COUNT(*) / 2
) FOO;
```

**(b) NULL を含まないグループだけ抽出**
```sql
SELECT dept FROM Personnel
GROUP BY dept
HAVING COUNT(*) = COUNT(salary);
-- 全行が NULL でない = 行数と非NULL値の数が一致
```

**(c) 関係除算（バスケット解析）**
「指定したアイテムを **すべて** 含むバスケットを探す」というのが関係除算。**HAVING + COUNT** か **NOT EXISTS の二重否定** で書く。
```sql
-- ItemsToBuy のアイテムを全部含む店舗
SELECT shop FROM ShopItems SI, ItemsToBuy I
WHERE SI.item = I.item
GROUP BY shop
HAVING COUNT(SI.item) = (SELECT COUNT(item) FROM ItemsToBuy);
```

**(d) すべての学生が条件を満たすか（全称量化）**
HAVING の中で `COUNT(*) = SUM(CASE WHEN cond THEN 1 ELSE 0 END)` のように書くと「全要素が cond を満たす」を表せる。

### コラム
- **関係除算**: 関係代数の「割り算」。商集合を求める操作で、HAVING の真骨頂。
- **HAVING句とウィンドウ関数の使い分け**: 集約して1行にしてよいなら HAVING、各行を残したいならウィンドウ関数 + 後置フィルタ。

### 学習のコツ
**ベン図を描く**。集合 A・B の和・積・差で表現できないかを先に考えるクセが、SQL を集合の言語として使うコツ。

---

## 7章 ウィンドウ関数で行間比較を行なう

### この章の主張
2章で導入したウィンドウ関数を、特に **「前の行・次の行・同パーティション内のN行前後」との比較** に活用する応用編。かつて自己結合や相関サブクエリで書いていた処理を、`LAG` / `LEAD` で劇的にシンプルにできる。

### 主な関数
- **LAG(col, n, default)**: n行前の値を取得（パーティション内 ORDER BY 基準）。
- **LEAD(col, n, default)**: n行後の値を取得。
- **FIRST_VALUE(col) / LAST_VALUE(col)**: ウィンドウの先頭/末尾の値。
- **NTH_VALUE(col, n)**: n番目の値。

### 典型パターン

**(a) 前日との売上差分**
```sql
SELECT order_date, amount,
  amount - LAG(amount) OVER (ORDER BY order_date) AS diff_from_prev
FROM DailySales;
```

**(b) 連続した期間の検出**
```sql
SELECT id, year, sales,
  CASE WHEN sales > LAG(sales) OVER (PARTITION BY id ORDER BY year)
       THEN '増加'
       WHEN sales < LAG(sales) OVER (PARTITION BY id ORDER BY year)
       THEN '減少'
       ELSE '横ばい'
  END AS trend
FROM YearlySales;
```

**(c) 累積比率（パーセンタイル風）**
```sql
SELECT name, score,
  CUME_DIST() OVER (ORDER BY score) AS cum_dist,
  PERCENT_RANK() OVER (ORDER BY score) AS pct_rank
FROM Scores;
```

### コツ
- **LAG / LEAD は ORDER BY 必須**（順序が決まらないと意味不明）。
- 日付に欠落がある場合、「1行前」と「1日前」は別物なので注意。日付ベースなら `RANGE BETWEEN INTERVAL '1' DAY PRECEDING ...` の方が安全。

---

## 8章 外部結合の使い方

### この章の主張
外部結合（`LEFT/RIGHT/FULL OUTER JOIN`）を「結合」と呼ぶより、**「行のレイアウトを変える整形ツール」** と捉える方が本質に近い。帳票・レポート作成で多用される。

### 典型パターン

**(a) 行 → 列（クロス集計表）**
```sql
SELECT M.name,
  T1.score AS Math,
  T2.score AS Science,
  T3.score AS English
FROM Members M
LEFT JOIN Scores T1 ON M.id = T1.id AND T1.subject = 'Math'
LEFT JOIN Scores T2 ON M.id = T2.id AND T2.subject = 'Science'
LEFT JOIN Scores T3 ON M.id = T3.id AND T3.subject = 'English';
```
1教科ごとに外部結合してメンバーに「列」をぶら下げる。CASE 式版（5章）と並ぶピボットの定石。

**(b) マスタを軸に「0件側」も出すレポート**
```sql
SELECT C.category_name, COUNT(P.product_id) AS num_products
FROM Categories C
LEFT JOIN Products P ON C.id = P.category_id
GROUP BY C.category_name;
```
INNER JOIN だと商品0個のカテゴリが消えてしまうが、LEFT JOIN なら 0 のまま残る。

**(c) 入れ子の外部結合で多段ピボット**
3つ以上のテーブルを外部結合で繋ぐ場合、結合順や ON 句の位置で結果が変わる。**括弧で結合順を明示**し、結合条件はそれぞれの ON に書くのが安全。

### よくある誤り
- LEFT JOIN したテーブルの列に `WHERE` で条件を付けると、**実質 INNER JOIN になる**。「LEFT JOIN 側の条件は ON に、ベース側の条件は WHERE に」が原則。

---

## 9章 SQLで集合演算

### この章の主張
SQL は本来「集合の言語」。`UNION` / `INTERSECT` / `EXCEPT (MINUS)` の集合演算をきちんと使えば、行レベルの処理に逃げずに**集合レベルで宣言的に**書ける。

### 演算の優先順位と注意
- 順位: `INTERSECT` > `UNION` ≒ `EXCEPT`（標準では）。明示的に括弧で囲むのが安全。
- **重複除去のコストは大きい**。重複が出ないと分かっているなら **`UNION ALL`** を使う（ソート省略で高速）。
- DBMS差: `INTERSECT`/`EXCEPT` は MySQL では非サポート（バージョン依存。8.0.31 以降で対応）。Oracle では `EXCEPT` の代わりに `MINUS`。

### 典型パターン

**(a) 差集合で「片方にしかない」レコード**
```sql
SELECT id FROM TableA
EXCEPT
SELECT id FROM TableB;
```

**(b) テーブルの相等性チェック**
2つのテーブル A, B が完全に等しいことを確かめる。
```sql
-- 空集合なら等しい
(SELECT * FROM A UNION SELECT * FROM B)
EXCEPT
(SELECT * FROM A INTERSECT SELECT * FROM B);
```
あるいは「行数が等しく、かつ和集合の行数も両者の行数と等しい」でもよい。

**(c) MySQL等で EXCEPT が使えないとき**
```sql
SELECT id FROM TableA A
WHERE NOT EXISTS (SELECT * FROM TableB B WHERE B.id = A.id);
```

### 演算結果のソート
集合演算の結果は順序が保証されない。並べたければ最後に `ORDER BY` を1回だけ書く。

---

## 10章 SQLで数列を扱う

### この章の主張
連番・欠番・連続区間といった「数列っぽい問題」を、ループや変数なしに**集合で解く**章。SQL の宣言型らしさが最もよく出るテーマの1つ。

### 典型パターン

**(a) 欠番の洗い出し（4通り）**
1〜100 のうち、テーブル SeqTbl に存在しない番号を出す：
```sql
-- ① NOT IN
SELECT n FROM Numbers WHERE n NOT IN (SELECT seq FROM SeqTbl);
-- ② NOT EXISTS
SELECT n FROM Numbers
WHERE NOT EXISTS (SELECT * FROM SeqTbl WHERE seq = n);
-- ③ EXCEPT
SELECT n FROM Numbers EXCEPT SELECT seq FROM SeqTbl;
-- ④ 外部結合
SELECT N.n FROM Numbers N
LEFT JOIN SeqTbl S ON N.n = S.seq
WHERE S.seq IS NULL;
```
**性能・NULL耐性・移植性で違いがある**。NULL耐性は ② > ④ > ③ > ①、移植性は ① ≒ ② が高い。

**(b) 連続範囲の検出（開始・終了）**
SeqTbl に 1,2,3, 5,6, 8,9,10 が入っているとき、`{1-3, 5-6, 8-10}` を出したい。
```sql
SELECT seq AS start_, MIN(S2.seq) AS end_
FROM SeqTbl S1, SeqTbl S2
WHERE S2.seq >= S1.seq
  AND NOT EXISTS (SELECT * FROM SeqTbl S3 WHERE S3.seq = S2.seq + 1)
  AND NOT EXISTS (SELECT * FROM SeqTbl S4 WHERE S4.seq = S1.seq - 1)
GROUP BY S1.seq;
```
「自分の前の番号がない = 開始点」「自分の次の番号がない = 終了点」を NOT EXISTS で表現する。

**(c) 全座席が空席の連続範囲を取れるか（全称量化）**
「3席連続して空いている範囲」を「空いていない座席が範囲内に1つも存在しない」と言い換えて NOT EXISTS。

### 学びどころ
ここまでくると「**ループっぽい問題を NOT EXISTS で言い換える**」型が身についてくる。これは19章の「7箇条」で総括される思想の応用編。

---

## 11章 SQLを速くするぞ

### この章の主張
**「正しい結果を出す SQL」は1段階目、「速い SQL」が2段階目**。インデックスを生かす書き方、コストの高い処理を減らす書き方を体系的に解説する実務必須章。

### インデックスを殺さない書き方
- **索引列に関数・演算を加えない**:
  - 悪: `WHERE SUBSTR(name, 1, 1) = 'A'`
  - 良: `WHERE name LIKE 'A%'`
  - 悪: `WHERE TRUNC(created_at) = '2024-01-01'`
  - 良: `WHERE created_at >= '2024-01-01' AND created_at < '2024-01-02'`
- **否定形は基本インデックス不可**: `<>`, `!=`, `NOT IN`, `IS NOT NULL`。
- **LIKE は前方一致のみ** インデックスが効く。`'%foo'` や `'%foo%'` はフルスキャン。
- **暗黙の型変換を避ける**: 文字列カラムに数値リテラルを書くと変換でインデックス無効化。

### 書き換えテクニック
- **`IN (サブクエリ)` → `EXISTS`**: 早期終了で速い場合が多い。
- **サブクエリ → 結合**: オプティマイザが処理しやすくなる。
- **`DISTINCT` → `EXISTS`**: 重複除去のためのソートを回避。
- **`UNION` → `UNION ALL`**（重複が出ないと分かるなら）。
- **`WHERE` で絞ってから `GROUP BY`**: グルーピング対象を減らす。
- **OR を UNION ALL に書き換える**: OR は索引利用が崩れることがある。

### 集約・順序付け
- 集計が頻繁に行われるならビュー / マテリアライズドビュー。
- `ORDER BY` のために `SORT` が発生する場合、インデックスの並びを利用できないか検討。

### 実行計画を読む
DB ごとの `EXPLAIN` / `EXPLAIN PLAN` / `SET STATISTICS IO` を使って**「想定通りインデックスが使われているか」を確認**。書き換えは仮説ベースではなく**計画ベース**で行う。

### 大原則
> **コストの高い処理（ソート・全表スキャン・大規模ハッシュ結合）を、より早い段階で減らす**。
> インデックス・絞り込み・結合順序のすべては、この目的のための道具立て。

---

## 12章 SQLプログラミング作法

### この章の主張
SQL も他の言語と同様、**可読性・保守性・移植性** を意識した書き方の作法がある。「動けばいい」を脱して、チームで運用できる SQL を書くための心得集。

### 命名・整形
- 列名は冗長でも構わないので**意味が読み取れる名前**を。略語を増やしすぎない。
- インデント・改行を揃え、句（SELECT / FROM / WHERE / GROUP BY ...）を左揃え。
- 1行が長くなったら **JOIN ごとに改行**、`AND` / `OR` で改行。

### 標準SQLを基本に、方言は局所化
- 移植性を保つために標準SQL機能で書けるところは書く。
- どうしても DB 固有機能を使うなら、**ビューやストアドプロシージャに閉じ込めて**呼び出し側を抽象化する。

### 安全な DML / DDL
- `DELETE` / `UPDATE` の **WHERE 漏れ事故** を避ける（トランザクション内で SELECT で確認してから実行する習慣）。
- 大量更新はチャンク分割を検討。
- DDL も冪等に書けるよう、`IF NOT EXISTS` などを活用。

### コメントと意図
- WHY を書く。WHAT はコードを読めば分かる。
- 「この WHERE 条件は◯◯バグの workaround」など、後から読む人が困る箇所には残す。

---

# 第2部 リレーショナルデータベースの世界

## 13章 RDB近現代史

### 何の章か
SQL を技術的に学ぶ前に、**なぜ RDB が主流になったか**、**今 RDB が直面している壁は何か**を歴史軸で押さえる。SQL の設計判断は、当時の技術的・思想的背景を抜きには理解できない。

### データモデルの系譜
- **階層型DB**（1960s, IMS）: 木構造。親→子のポインタ参照が中心。
- **ネットワーク型DB**（CODASYL, 1970s）: 多対多が表現できる代わりに、参照経路が複雑。
- **リレーショナルDB**（1970, コッドが提唱）: 表形式・SQL で簡潔に書ける。**第1の破壊的イノベーション**。

RDB が勝った理由：直感的なデータモデル（表）+ ユーザフレンドリーな宣言型言語（SQL）。

### 現代の課題
1. **性能と信頼性のトレードオフ**: ACID を強く守るほど、ストレージ・ロックが**シングルポイント**になり、スケールアウトが効きにくい。
2. **データモデルの限界**: 表形式は**グラフ構造・半構造化データ**（JSON, ドキュメント、画像）には不向き。

### NoSQL の位置づけ
- データモデルを単純化（KVS, ドキュメント, グラフ）。
- ACID を緩めて（多くは BASE: Basically Available, Soft state, Eventually consistent）スケールアウト性を獲得。
- 本書スタンス：**「第2の破壊的イノベーションの候補」だが、まだ RDB を完全に置き換えるには至っていない**。

---

## 14章 なぜ"関係"モデルという名前なの?

### この章の主張
**「表」と「関係（リレーション）」は同じものではない**。RDB が「リレーショナル」と呼ばれるのは、その背後に**集合論の関係（Cartesian Product の部分集合）**という数学的基礎があるから。

### 関係の性質（表とは違う点）
- **行に順序がない**: SQL の `ORDER BY` がない限り行の順序は保証されない。
- **列（属性）にも順序がない**: 関係モデル上は属性の並びは意味を持たない。
- **重複行が存在しない**: 関係は集合だから。SQL の表は**マルチセット**（重複OK）なので、ここに大きなズレがある。
- **値はすべて単一値**（第1正規形）: 配列やネストは認めない（純粋関係モデルの立場）。

### モデルと実装のズレ
SQL は関係モデルの**緩い実装**であり、純粋モデルとの差を知っておくことが大事。たとえば SQL では `SELECT * FROM Tbl` で重複が出ることがある（NULL も絡む）。これがバグの温床になる。

---

## 15章 関係に始まり関係に終わる

### この章の主張
**関係代数の閉包性**: 関係に演算（選択・射影・結合・和・差…）を施した結果も**また関係**である。これにより、結果をさらに演算の入力にする「合成」が自然にできる。

### SQL での具体化
- `SELECT` の結果も**表（≒関係）**。だから `SELECT ... FROM (SELECT ...) AS sub` のように**サブクエリで合成**できる。
- ビュー、CTE (`WITH`)、派生表 (`FROM (SELECT ...)`) は全てこの閉包性の応用。

### アナロジー
UNIX パイプ `cat foo | grep bar | sort` も、各コマンドが**ファイル → ファイル**で閉じているから合成できる。SQL も「集合 → 集合」で閉じているから合成できる。**関数型プログラミングの合成（compose）と同じ発想**。

---

## 16章 アドレス、この巨大な怪物

### この章の主張
RDB の最大の発明は、データから **アドレス（ポインタ）を追放したこと**。コッドの目的は「**ループを書かなくて済むデータベースを作ること**」だった。

### ポインタの何が問題か
- ポインタを辿るには**順序づけられた処理（ループ）**が必要 → 手続き型に逆戻り。
- ポインタは**物理レイアウト**に依存 → データ移動・再編成のたびに壊れる。
- 並列実行・最適化を阻害する。

### RDB の選択
- レコード同士の関係は**値（外部キー）** で表す。
- アクセスは**集合演算**で行う。物理アクセス経路はオプティマイザに任せる。
- 結果として、**プログラマは「何を取りたいか（What）」だけ書けばよく、「どう取るか（How）」は DB に委ねられる**。

### 関数型言語との親和性
Lisp / Haskell / Scala のような関数型言語は、**副作用なし・参照透過・ループより再帰や写像**という思想で、SQL とよく似ている。手続き型の感覚で SQL を書こうとすると詰まるのは、根本的にパラダイムが違うから。

---

## 17章 順序をめぐる冒険

### この章の主張
関係モデルでは**行に順序はない**。なのに実務では「**前年比、移動平均、累積、ランキング**」と順序を扱いたい場面が山ほどある。この**矛盾をどう折り合わせるか**が、ウィンドウ関数登場の歴史的背景。

### 矛盾の解消
- SQL の基本構造（GROUP BY, 結合）は順序を扱わない。
- OLAP（分析処理）の需要に応えるため、SQL:1999 以降で**ウィンドウ関数**が標準化された。
- ウィンドウ関数は「**順序を一時的に持ち込んだ局所領域**」で集約を行う。関係モデルの純粋性を守りつつ、実務の要求を満たす設計。

### 設計上の示唆
- テーブル設計時には**「順序を持つ意味のある列（時系列、版番号）」を必ず明示**する。`id`連番に依存して順序を仮定するのは危険。
- 「順序が要るかどうか」を意識して、必要なら ORDER BY や ウィンドウ関数の OVER (ORDER BY ...) を明示。

---

## 18章 GROUP BYとPARTITION BY

### この章の主張
**両者は似て非なるもの**。違いを一言で言えば、**潰すか・潰さないか**。

| | GROUP BY | PARTITION BY |
|---|---|---|
| 入力 | 集合 | 集合 |
| 分割 | グループに分ける | パーティションに分ける |
| 集約 | **1グループ → 1行**に圧縮 | **圧縮しない**（行数は変わらない） |
| 用途 | 集約結果が欲しい | 行を残したまま集計値を付加 |

### 使い分けの例
- 「部署別の人数」 → GROUP BY で OK。
- 「各従業員の隣に、自分が所属する部署の人数を表示したい」 → PARTITION BY。

```sql
-- 各従業員と、その部署の人数
SELECT emp_id, dept,
  COUNT(*) OVER (PARTITION BY dept) AS dept_size
FROM Employees;
```

### 派生的理解
- GROUP BY は集合演算の延長。
- PARTITION BY は **「関係を保ったまま局所集計する」** という発想で、関係モデルに新しい次元を加える。

---

## 19章 手続き型から宣言型・集合指向へ頭を切り替える7箇条

### この章の主張
SQL を「読めるけど書けない」「書けるけど性能が出ない」状態から抜けるには、**思考のパラダイムシフト**が必要。本書全体のエッセンスを7つに凝縮した章。

### 7箇条
1. **テーブルの行に順序はない** — 順序が要るなら明示的に ORDER BY / ウィンドウ関数。
2. **ループは GROUP BY とウィンドウ関数で置き換える** — for/while 思考を捨てる。
3. **EXISTS と「量化」を理解する** — ∃/∀ を意識して述語論理として読む。
4. **HAVING句が SQL の集合指向の本質** — 集合の性質を問う武器。
5. **IF文は CASE式で置き換える** — 式として組み合わせる。
6. **テーブルを行の列挙ではなく「集合」とみなす** — 操作対象は集合。
7. **ベン図を描いて考える** — 和・積・差で表現できないかをまず考える。

### なぜパラダイムシフトが重い?
新人研修でも、**Python や Java よりも SQL の方が「脳への負荷」が高い** と言われる。多くの人が手続き型を先に学ぶため、宣言型・集合指向への切り替えに認知コストがかかる。

---

## 20章 神のいない論理

### この章の主張
古典論理（アリストテレス以来）は **真/偽の2値** ＝ 「**全知の神の論理**」。だが現実世界には「分からない」が常に存在する。SQL の **3値論理は「不完全情報を扱うための、神のいない論理」**。

### 3値論理の必然性
- 現実のデータは不完全（未入力、未測定、適用不能）。
- それらを「偽」と同一視するのは情報の捏造。
- だから「**UNKNOWN（不明）**」という第3の真偽値が要る。

### 代償
- 直感に反する挙動が増える（4章で見た NOT IN の罠など）。
- 開発者が常に「UNKNOWN を想定する」必要があり、思考コストが高い。

### 教訓
3値論理は**便利だが危険な道具**。**設計時に NULL を増やさないことで、結局は3値論理を回避するのが現実解**。

---

## 21章 SQLと再帰集合

### この章の主張
SQL は **再帰共通表式（`WITH RECURSIVE`）** で木構造・グラフ探索を扱える。これは集合論的な**再帰的集合の定義**にきちんと裏打ちされている。

### 数学的下敷き
- フォン・ノイマンによる自然数の構成:
  - 0 := ∅（空集合）
  - n + 1 := n ∪ { n }
- 「初期値 + 漸化式」で集合を構築する発想は、`WITH RECURSIVE` の構文そのもの:
  ```sql
  WITH RECURSIVE Tree AS (
    SELECT id, parent_id, 1 AS depth FROM Nodes WHERE parent_id IS NULL
    UNION ALL
    SELECT N.id, N.parent_id, T.depth + 1
    FROM Nodes N JOIN Tree T ON N.parent_id = T.id
  )
  SELECT * FROM Tree;
  ```

### 実務応用
- 組織図、部品表（BOM）、カテゴリ階層、SNSフォロー関係などのグラフ探索。
- ループしないように `UNION` で重複排除、または訪問済みフラグを持ち回る設計が必要。

---

## 22章 NULL撲滅委員会

### この章の主張
> **「NULL は薬。正しく使えば有用だが、乱用すれば全てをぶち壊す」**

NULL を**徹底的に減らす**設計が結局はバグも性能問題も減らす、という強い主張をまとめた章。

### NULL の害（おさらい）
- **3値論理の伝染**: NULL が1個入るだけで、関連クエリの結果が信頼できなくなる。
- **インデックスの利きにくさ**: 多くの DB で NULL は索引に含まれない、または挙動が特殊。
- **集約関数の例外的挙動**: COUNT(*) と COUNT(col) の差で集計値がずれる。
- **`NOT IN` で結果が消滅** などの罠。

### NULL を減らす設計指針
1. **NOT NULL 制約を基本に据える**。
2. デフォルト値を設定する（不明なら `'UNKNOWN'`, 未定なら `9999-12-31` など）。
3. **「未知」と「適用不能」を区別したいなら、別テーブルかフラグで表現**。たとえば「退職者の最終勤務日」を NULL にする代わりに、現役/退職者テーブルを分ける。
4. 集計対象列はとくに NULL 禁止が望ましい。

### NULL がどうしても必要なとき
- 値域が明確に「不明」を意味するケース（センサーデータの欠測など）。
- そのときも **`COALESCE` / `NVL` / `IFNULL` で常に NULL を変換**してからクエリを書く習慣。

---

## 23章 SQLにおける存在の階層

### この章の主張
SQL は「ある / ない」の表現が**多層化**している。それぞれの違いを意識して使い分けないと、思わぬバグになる。

### 存在の階層（ざっくり）
| レベル | 意味 | 表現 |
|---|---|---|
| 行が存在しない | テーブルに該当行がない | `NOT EXISTS`, 外部結合の NULL 行 |
| 集合として空 | サブクエリの結果が空集合 | `COUNT(*) = 0`, `NOT EXISTS` |
| 値として NULL | 行はあるが値が不明 | `IS NULL` |
| 空文字列 | 行も値もあるが、内容が空 | `= ''` |
| ゼロ・偽値 | 値はあるが意味的に「無」 | `= 0`, `= FALSE` |

### 落とし穴
- **「行が無い」と「値が NULL」を混同**: 外部結合の結果の NULL は「マッチする行がなかった」のサインで、元データの NULL とは別物。
- **空集合の集約**: `SELECT SUM(col) FROM Empty` は NULL を返す（0 ではない）。`COALESCE` で 0 に変換するのが安全。

### この章の意味
最終章として「**SQL は集合・述語論理・3値論理が絡み合った「存在を扱うための言語」**」というメッセージで本書を締めくくる。

---

# 第3部 付録

- **A. 演習問題の解答** — 各章末問題の解答例と解説。
- **B. 参考文献** — Date, Celko, Codd ら関係モデル理論の主要文献、SQL 標準仕様、関連書籍へのポインタ。

---

# 全体総括

### 本書の核となるメッセージ
- SQL は**手続き型の浅い理解では永遠に中級者になれない**。
- **関係モデル・集合論・述語論理** が裏にあると意識しながら書くと、コードが急に短く・速く・正しくなる。
- **ウィンドウ関数** は SQL の表現力を一段引き上げた革命的機能。第2版を貫くテーマ。
- **NULL は減らす**。

### 学習導線の提案
1. まず **1章 CASE式** と **2章 ウィンドウ関数** で「式」と「集合に対する集計」の感覚を掴む。
2. **4章 NULL** と **5章 EXISTS** で論理の罠を回避できるようになる。
3. **6章 HAVING** と **9章 集合演算** で集合指向の思考に慣れる。
4. **11章 チューニング** を実務でループ参照する。
5. 第2部は理論編。腹落ちさせたい時に何度も読む。

### 本書のあとに読むと良いもの
- 同著者『SQLパズル 第2版』(Celko 翻訳) — 演習中心で集合指向SQL力を鍛える。
- 同著者『達人に学ぶDB設計徹底指南書』 — 設計編。本書は SQL 編なので両輪。
- Joe Celko 『プログラマのためのSQL 第4版』 — 集合指向SQLのバイブル。

---

# 情報源（このまとめの参考）

- 翔泳社 書誌ページ: https://www.shoeisha.co.jp/book/detail/9784798157825
- CodeZine 紹介記事: https://codezine.jp/article/detail/11089
- Qiita 技術書まとめ (shimpeitakeda55): https://qiita.com/shimpeitakeda55/items/c6d9def1f8f3d940b8f2
- Zenn 読了記事 (churadata): https://zenn.dev/churadata/articles/9bbc25cf0127b6
- スマーティブ 書籍紹介: https://smartive.co.jp/book/262
- 著者ミック氏 Web連載「CASE式のススメ」: https://mickindex.sakura.ne.jp/database/db_case.html
- 著者ミック氏 サポートページ: https://mickindex.sakura.ne.jp/database/db_support_sinan.html
- Chapter3 自己結合 読書メモ: https://wand-ta.hatenablog.com/entry/2018/12/02/232551
- Qiita 3値論理: https://qiita.com/ryosuketter/items/e9eadaeb48540e94d88b
- Qiita 無料で学ぶ第1版 (連載リンク集): https://qiita.com/katayamahide/items/48f7a78dab3497adcd0c

# 本リポジトリ内の関連メモ

- [[tatsujin-db-design-toc]] — 姉妹書『達人に学ぶDB設計徹底指南書』目次
- [[db-design-tatsujin-mapping]] — 同上 マッピング
- [[db-refactoring-by-tatsujin-chapter]] — 同上 章別整理
