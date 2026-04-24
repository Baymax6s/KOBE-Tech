DELETE FROM articles
WHERE title IN (
    '神戸大学でのハッカソン体験記',
    'Goで作るREST API入門',
    'Vue 3 + Vuetifyで学ぶフロントエンド開発',
    'PostgreSQLのマイグレーション管理',
    'Dockerで開発環境を統一する'
);
