INSERT INTO articles (title, content, user_id, created_at, updated_at)
SELECT v.title, v.content, u.id, v.created_at, v.updated_at
FROM (VALUES
    ('神戸大学でのハッカソン体験記',
     '先日、神戸大学で開催されたハッカソンに参加しました。チームメンバーと48時間かけてWebアプリを開発し、多くのことを学びました。特にチーム開発の難しさと楽しさを実感できた貴重な経験でした。',
     'admin',
     '2026-04-01 09:00:00+09'::timestamptz,
     '2026-04-01 09:00:00+09'::timestamptz),
    ('Goで作るREST API入門',
     'GoのginフレームワークでREST APIを作る基本を紹介します。ルーティング、ミドルウェア、JSONレスポンスの返し方など、実際のコードを交えながら解説します。Goは静的型付けと高いパフォーマンスが魅力で、バックエンド開発に最適です。',
     'user01',
     '2026-04-10 12:30:00+09'::timestamptz,
     '2026-04-10 12:30:00+09'::timestamptz),
    ('Vue 3 + Vuetifyで学ぶフロントエンド開発',
     'Vue 3のComposition APIとVuetifyを組み合わせたフロントエンド開発の入門記事です。コンポーネント設計やリアクティブな状態管理の基本を、サンプルコードを通じて説明します。',
     'user02',
     '2026-04-15 15:00:00+09'::timestamptz,
     '2026-04-15 15:00:00+09'::timestamptz),
    ('PostgreSQLのマイグレーション管理',
     'golang-migrateを使ったDBマイグレーションの管理方法を解説します。up/downファイルの書き方やCIへの組み込み方など、チーム開発で役立つプラクティスを紹介します。',
     'user03',
     '2026-04-18 10:00:00+09'::timestamptz,
     '2026-04-18 10:00:00+09'::timestamptz),
    ('Dockerで開発環境を統一する',
     'docker composeを使って開発環境を統一する方法を紹介します。PostgreSQLのコンテナ起動やボリューム管理など、チームで環境差異をなくすためのTipsをまとめました。',
     'user01',
     '2026-04-20 11:00:00+09'::timestamptz,
     '2026-04-20 11:00:00+09'::timestamptz)
) AS v(title, content, user_name, created_at, updated_at)
JOIN users u ON u.name = v.user_name
WHERE NOT EXISTS (
    SELECT 1 FROM articles a WHERE a.title = v.title AND a.user_id = u.id
);
