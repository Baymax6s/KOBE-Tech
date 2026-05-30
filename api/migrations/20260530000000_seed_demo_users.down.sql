-- 本マイグレーションで追加したデモユーザーとそのプロフィールのみを削除する。
-- 投稿・返信・いいねは後続マイグレーション (010000/020000/030000) の down で先に消える前提。
DELETE FROM user_profiles
WHERE user_id IN (
    SELECT id FROM users
    WHERE name IN ('鈴木健一', '高橋美咲', '渡辺翔太', '伊藤さくら', '中村大輔', '小林優子')
);

DELETE FROM users
WHERE name IN ('鈴木健一', '高橋美咲', '渡辺翔太', '伊藤さくら', '中村大輔', '小林優子');
