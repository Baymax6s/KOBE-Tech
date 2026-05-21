-- 質問 / 回答 (kind = 1, 2) を追加投入する。
-- 既存の seed_replies では Article 1, 2 に少量の Q&A しかなかったため、
-- 残り 3 記事と既存記事にもネスト構造を含む Q&A を増やし、
-- フロントの「質問/回答」フローと折りたたみ UI を実データで確認できるようにする。
DO $$
DECLARE
    v_article1_id INT; -- 神戸大学でのハッカソン体験記
    v_article2_id INT; -- Goで作るREST API入門
    v_article3_id INT; -- Vue 3 + Vuetifyで学ぶフロントエンド開発
    v_article4_id INT; -- PostgreSQLのマイグレーション管理
    v_article5_id INT; -- Dockerで開発環境を統一する

    v_admin_id    INT;
    v_user01_id   INT;
    v_user02_id   INT;
    v_user03_id   INT;

    v_a1_q_id   INT;
    v_a1_a1_id  INT;
    v_a3_q_id   INT;
    v_a3_a1_id  INT;
    v_a4_q_id   INT;
    v_a4_a1_id  INT;
    v_a5_q_id   INT;
    v_a5_a1_id  INT;
BEGIN
    -- 冪等性: 本マイグレーション特有の質問本文の存在で判定する
    IF EXISTS (
        SELECT 1 FROM replies
        WHERE content = 'Vuetify の v-data-table をカスタマイズしたいのですが、まず何から触るのがおすすめですか？'
    ) THEN
        RETURN;
    END IF;

    SELECT id INTO v_article1_id FROM articles WHERE title = '神戸大学でのハッカソン体験記';
    SELECT id INTO v_article2_id FROM articles WHERE title = 'Goで作るREST API入門';
    SELECT id INTO v_article3_id FROM articles WHERE title = 'Vue 3 + Vuetifyで学ぶフロントエンド開発';
    SELECT id INTO v_article4_id FROM articles WHERE title = 'PostgreSQLのマイグレーション管理';
    SELECT id INTO v_article5_id FROM articles WHERE title = 'Dockerで開発環境を統一する';

    SELECT id INTO v_admin_id  FROM users WHERE name = 'admin';
    SELECT id INTO v_user01_id FROM users WHERE name = '田中太郎';
    SELECT id INTO v_user02_id FROM users WHERE name = '山田花子';
    SELECT id INTO v_user03_id FROM users WHERE name = '佐藤次郎';

    -- ========== Article 1: 神戸大学でのハッカソン体験記 ==========
    -- 質問 → 回答 → さらに回答（追問）
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_user02_id,
            '次回参加するときに、当日までに練習しておくと良い技術スタックはありますか？',
            1, NULL,
            '2026-04-05 10:00:00+09'::timestamptz,
            '2026-04-05 10:00:00+09'::timestamptz)
    RETURNING id INTO v_a1_q_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_admin_id,
            'チーム構成にもよりますが、最低限 Git の運用と、フロント・バックの最小サンプルを動かせるようにしておくと当日が楽です。',
            2, v_a1_q_id,
            '2026-04-05 12:00:00+09'::timestamptz,
            '2026-04-05 12:00:00+09'::timestamptz)
    RETURNING id INTO v_a1_a1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_user02_id,
            'なるほど、フロントとバックの最小サンプルというのは、具体的にはどんな構成でしょうか？',
            2, v_a1_a1_id,
            '2026-04-05 13:30:00+09'::timestamptz,
            '2026-04-05 13:30:00+09'::timestamptz);

    -- ========== Article 3: Vue 3 + Vuetifyで学ぶフロントエンド開発 ==========
    -- 質問 → 回答（ベストアンサー） + 別の回答
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article3_id, v_user01_id,
            'Vuetify の v-data-table をカスタマイズしたいのですが、まず何から触るのがおすすめですか？',
            1, NULL,
            '2026-04-17 10:00:00+09'::timestamptz,
            '2026-04-17 10:00:00+09'::timestamptz)
    RETURNING id INTO v_a3_q_id;

    INSERT INTO replies (article_id, user_id, content, is_best, kind, parent_id, created_at, updated_at)
    VALUES (v_article3_id, v_admin_id,
            'まずは headers の定義と #item.xxx スロットでの上書きから始めると、見た目を壊さずに置き換えられます。',
            TRUE, 2, v_a3_q_id,
            '2026-04-17 12:00:00+09'::timestamptz,
            '2026-04-17 12:00:00+09'::timestamptz)
    RETURNING id INTO v_a3_a1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article3_id, v_user02_id,
            '#item.xxx の他にも、行全体を差し替えたい場合は #item スロットを使うこともできますよ。',
            2, v_a3_q_id,
            '2026-04-17 14:00:00+09'::timestamptz,
            '2026-04-17 14:00:00+09'::timestamptz);

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article3_id, v_user01_id,
            'スロットを段階的に試すと挙動が理解しやすいんですね、参考にします！',
            2, v_a3_a1_id,
            '2026-04-17 15:30:00+09'::timestamptz,
            '2026-04-17 15:30:00+09'::timestamptz);

    -- ========== Article 4: PostgreSQLのマイグレーション管理 ==========
    -- 質問 → 回答（ベストアンサー）
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article4_id, v_user03_id,
            '本番でマイグレーションを流すタイミングは、デプロイ前と後どちらにするのが良いでしょうか？',
            1, NULL,
            '2026-04-21 09:00:00+09'::timestamptz,
            '2026-04-21 09:00:00+09'::timestamptz)
    RETURNING id INTO v_a4_q_id;

    INSERT INTO replies (article_id, user_id, content, is_best, kind, parent_id, created_at, updated_at)
    VALUES (v_article4_id, v_admin_id,
            '追加系のスキーマ変更ならデプロイ前、列削除など破壊的なものは新コードのデプロイ後に流す二段階運用が安全です。',
            TRUE, 2, v_a4_q_id,
            '2026-04-21 11:00:00+09'::timestamptz,
            '2026-04-21 11:00:00+09'::timestamptz)
    RETURNING id INTO v_a4_a1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article4_id, v_user03_id,
            '破壊的変更を後ろに回すのは盲点でした。実運用で気をつけます。',
            2, v_a4_a1_id,
            '2026-04-21 12:30:00+09'::timestamptz,
            '2026-04-21 12:30:00+09'::timestamptz);

    -- ========== Article 5: Dockerで開発環境を統一する ==========
    -- 質問 → 回答 → 回答（議論っぽいネスト）
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article5_id, v_user01_id,
            'M1 Mac だと一部のイメージで platform 警告が出るのですが、皆さんどう対処していますか？',
            1, NULL,
            '2026-04-25 10:00:00+09'::timestamptz,
            '2026-04-25 10:00:00+09'::timestamptz)
    RETURNING id INTO v_a5_q_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article5_id, v_user02_id,
            'compose で platform: linux/amd64 を明示しています。ビルドは遅くなりますが、本番との差分を抑えられます。',
            2, v_a5_q_id,
            '2026-04-25 11:30:00+09'::timestamptz,
            '2026-04-25 11:30:00+09'::timestamptz)
    RETURNING id INTO v_a5_a1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article5_id, v_admin_id,
            '私はマルチアーキ対応の公式イメージを優先して選ぶ運用にしていて、できるだけ platform 指定を避けています。',
            2, v_a5_q_id,
            '2026-04-25 13:00:00+09'::timestamptz,
            '2026-04-25 13:00:00+09'::timestamptz);

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article5_id, v_user01_id,
            'なるほど、両方試して開発時の体感とのバランスを見てみます。ありがとうございます！',
            2, v_a5_a1_id,
            '2026-04-25 14:30:00+09'::timestamptz,
            '2026-04-25 14:30:00+09'::timestamptz);
END $$;
