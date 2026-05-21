-- 上記 seed で投入した (article, user) の組み合わせのみを削除する。
-- ユーザーが手動で押した本物のいいねを巻き込まないよう、組み合わせを明示する。
DELETE FROM likes l
USING articles a, users u
WHERE l.article_id = a.id
  AND l.user_id    = u.id
  AND (a.title, u.name) IN (
    ('神戸大学でのハッカソン体験記',          'user01'),
    ('神戸大学でのハッカソン体験記',          'user02'),
    ('神戸大学でのハッカソン体験記',          'user03'),
    ('Goで作るREST API入門',                  'admin'),
    ('Goで作るREST API入門',                  'user02'),
    ('Goで作るREST API入門',                  'user03'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'admin'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'user01'),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発', 'user03'),
    ('PostgreSQLのマイグレーション管理',      'admin'),
    ('PostgreSQLのマイグレーション管理',      'user01'),
    ('Dockerで開発環境を統一する',            'admin'),
    ('Dockerで開発環境を統一する',            'user02'),
    ('Dockerで開発環境を統一する',            'user03')
  );
