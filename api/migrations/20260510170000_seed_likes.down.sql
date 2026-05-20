-- 上記 seed で投入した (article, user) の組み合わせのみを削除する。
-- ユーザーが手動で押した本物のいいねを巻き込まないよう、組み合わせを明示する。
DELETE FROM likes l
USING articles a, users u
WHERE l.article_id = a.id
  AND l.user_id    = u.id
  AND (a.title, u.name) IN (
    ('神戸大学でのハッカソン体験記',          '田中太郎'),
    ('神戸大学でのハッカソン体験記',          '山田花子'),
    ('神戸大学でのハッカソン体験記',          '佐藤次郎'),
    ('Goで作るREST API入門',                  'admin'),
    ('Goで作るREST API入門',                  '山田花子'),
    ('Goで作るREST API入門',                  '佐藤次郎'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'admin'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', '田中太郎'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', '佐藤次郎'),
    ('PostgreSQLのマイグレーション管理',      'admin'),
    ('PostgreSQLのマイグレーション管理',      '田中太郎'),
    ('Dockerで開発環境を統一する',            'admin'),
    ('Dockerで開発環境を統一する',            '山田花子'),
    ('Dockerで開発環境を統一する',            '佐藤次郎')
  );
