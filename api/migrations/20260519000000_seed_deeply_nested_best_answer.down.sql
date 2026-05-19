-- 本マイグレーションで投入した回答のみを削除する。
DELETE FROM replies
WHERE content = '最小サンプルは「フロントが Hello を表示」「バックが GET /ping で pong を返す」「両者を fetch で繋ぐ」の 3 点セットで十分です。当日の認識合わせがぐっと早くなります。';
