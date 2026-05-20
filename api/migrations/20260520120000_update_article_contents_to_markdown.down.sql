-- Markdown 化前の plain text content に戻す。
UPDATE articles SET content = '先日、神戸大学で開催されたハッカソンに参加しました。チームメンバーと48時間かけてWebアプリを開発し、多くのことを学びました。特にチーム開発の難しさと楽しさを実感できた貴重な経験でした。'
WHERE title = '神戸大学でのハッカソン体験記';

UPDATE articles SET content = 'GoのginフレームワークでREST APIを作る基本を紹介します。ルーティング、ミドルウェア、JSONレスポンスの返し方など、実際のコードを交えながら解説します。Goは静的型付けと高いパフォーマンスが魅力で、バックエンド開発に最適です。'
WHERE title = 'Goで作るREST API入門';

UPDATE articles SET content = 'Vue 3のComposition APIとVuetifyを組み合わせたフロントエンド開発の入門記事です。コンポーネント設計やリアクティブな状態管理の基本を、サンプルコードを通じて説明します。'
WHERE title = 'Vue 3 + Vuetifyで学ぶフロントエンド開発';

UPDATE articles SET content = 'golang-migrateを使ったDBマイグレーションの管理方法を解説します。up/downファイルの書き方やCIへの組み込み方など、チーム開発で役立つプラクティスを紹介します。'
WHERE title = 'PostgreSQLのマイグレーション管理';

UPDATE articles SET content = 'docker composeを使って開発環境を統一する方法を紹介します。PostgreSQLのコンテナ起動やボリューム管理など、チームで環境差異をなくすためのTipsをまとめました。'
WHERE title = 'Dockerで開発環境を統一する';
