-- replies (コメント / 質問 / 回答) の seed
-- kind: 0 = comment, 1 = question, 2 = answer
-- parent_id を持つ返信は親の id を受ける必要があるため DO ブロックで記述する。
DO $$
DECLARE
    v_article1_id INT;
    v_article2_id INT;
    v_article3_id INT;
    v_admin_id    INT;
    v_user01_id   INT;
    v_user02_id   INT;
    v_user03_id   INT;
    v_a1_comment1_id  INT;
    v_a1_question1_id INT;
    v_a2_question1_id INT;
    v_a2_answer1_id   INT;
BEGIN
    -- 既に seed 済みなら何もしない（冪等性）
    IF EXISTS (
        SELECT 1
        FROM replies r
        JOIN articles a ON a.id = r.article_id
        WHERE a.title = '神戸大学でのハッカソン体験記'
          AND r.parent_id IS NULL
    ) THEN
        RETURN;
    END IF;

    SELECT id INTO v_article1_id FROM articles WHERE title = '神戸大学でのハッカソン体験記';
    SELECT id INTO v_article2_id FROM articles WHERE title = 'Goで作るREST API入門';
    SELECT id INTO v_article3_id FROM articles WHERE title = 'Vue 3 + Vuetifyで学ぶフロントエンド開発';
    SELECT id INTO v_admin_id  FROM users WHERE name = 'admin';
    SELECT id INTO v_user01_id FROM users WHERE name = '田中太郎';
    SELECT id INTO v_user02_id FROM users WHERE name = '山田花子';
    SELECT id INTO v_user03_id FROM users WHERE name = '佐藤次郎';

    -- ========== Article 1: 神戸大学でのハッカソン体験記 ==========
    -- コメント（記事直下） + 子コメント
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_user01_id,
            '楽しそうですね！どんなお題のハッカソンだったんですか？',
            0, NULL,
            '2026-04-02 10:00:00+09'::timestamptz,
            '2026-04-02 10:00:00+09'::timestamptz)
    RETURNING id INTO v_a1_comment1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_admin_id,
            '教育系のWebアプリ開発がテーマでした。発表もかなり盛り上がりましたよ！',
            0, v_a1_comment1_id,
            '2026-04-02 12:30:00+09'::timestamptz,
            '2026-04-02 12:30:00+09'::timestamptz);

    -- 質問（記事直下） + ベストアンサー
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_user02_id,
            'チームメンバーは何人くらいでしたか？役割分担も気になります。',
            1, NULL,
            '2026-04-03 09:00:00+09'::timestamptz,
            '2026-04-03 09:00:00+09'::timestamptz)
    RETURNING id INTO v_a1_question1_id;

    INSERT INTO replies (article_id, user_id, content, is_best, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_admin_id,
            '4人チームで、フロント2名・バック1名・デザイン1名の構成でした。',
            TRUE, 2, v_a1_question1_id,
            '2026-04-03 11:00:00+09'::timestamptz,
            '2026-04-03 11:00:00+09'::timestamptz);

    -- ========== Article 2: Goで作るREST API入門 ==========
    -- コメント（単体）
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article2_id, v_user02_id,
            '丁寧な記事ありがとうございます。手元でも動かせました！',
            0, NULL,
            '2026-04-11 10:00:00+09'::timestamptz,
            '2026-04-11 10:00:00+09'::timestamptz);

    -- 質問 → 回答 → 回答（ネスト）
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article2_id, v_user03_id,
            'gin と echo だと、最初に触るならどちらがおすすめですか？',
            1, NULL,
            '2026-04-12 14:00:00+09'::timestamptz,
            '2026-04-12 14:00:00+09'::timestamptz)
    RETURNING id INTO v_a2_question1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article2_id, v_user01_id,
            '用途次第ですが、最初はシンプルさで gin から入るのが追いやすいと思います。',
            2, v_a2_question1_id,
            '2026-04-12 16:00:00+09'::timestamptz,
            '2026-04-12 16:00:00+09'::timestamptz)
    RETURNING id INTO v_a2_answer1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article2_id, v_user03_id,
            'なるほど、ありがとうございます！週末に触ってみます。',
            2, v_a2_answer1_id,
            '2026-04-12 18:00:00+09'::timestamptz,
            '2026-04-12 18:00:00+09'::timestamptz);

    -- ========== Article 3: Vue 3 + Vuetifyで学ぶフロントエンド開発 ==========
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article3_id, v_admin_id,
            'Composition API のあたり、初学者にも追える説明でとても助かりました。',
            0, NULL,
            '2026-04-16 11:00:00+09'::timestamptz,
            '2026-04-16 11:00:00+09'::timestamptz);
END $$;
