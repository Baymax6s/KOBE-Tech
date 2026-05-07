UPDATE articles a
SET content = '先日、神戸大学で開催されたハッカソンに参加しました。チームメンバーと48時間かけてWebアプリを開発し、多くのことを学びました。特にチーム開発の難しさと楽しさを実感できた貴重な経験でした。'
FROM users u
WHERE a.user_id = u.id
  AND a.title = '神戸大学でのハッカソン体験記'
  AND u.name = 'admin';

UPDATE articles a
SET content = 'GoのginフレームワークでREST APIを作る基本を紹介します。ルーティング、ミドルウェア、JSONレスポンスの返し方など、実際のコードを交えながら解説します。Goは静的型付けと高いパフォーマンスが魅力で、バックエンド開発に最適です。'
FROM users u
WHERE a.user_id = u.id
  AND a.title = 'Goで作るREST API入門'
  AND u.name = 'user01';

UPDATE articles a
SET content = 'Vue 3のComposition APIとVuetifyを組み合わせたフロントエンド開発の入門記事です。コンポーネント設計やリアクティブな状態管理の基本を、サンプルコードを通じて説明します。'
FROM users u
WHERE a.user_id = u.id
  AND a.title = 'Vue 3 + Vuetifyで学ぶフロントエンド開発'
  AND u.name = 'user02';

UPDATE articles a
SET content = 'golang-migrateを使ったDBマイグレーションの管理方法を解説します。up/downファイルの書き方やCIへの組み込み方など、チーム開発で役立つプラクティスを紹介します。'
FROM users u
WHERE a.user_id = u.id
  AND a.title = 'PostgreSQLのマイグレーション管理'
  AND u.name = 'user03';

UPDATE articles a
SET content = 'docker composeを使って開発環境を統一する方法を紹介します。PostgreSQLのコンテナ起動やボリューム管理など、チームで環境差異をなくすためのTipsをまとめました。'
FROM users u
WHERE a.user_id = u.id
  AND a.title = 'Dockerで開発環境を統一する'
  AND u.name = 'user01';
