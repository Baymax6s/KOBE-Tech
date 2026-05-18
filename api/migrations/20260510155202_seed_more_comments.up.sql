-- 既存 seed に追加で「コメント」(kind = 'comment') を投入する。
-- 質問 / 回答 (kind = 'question', 'answer') は今回追加しない。
-- 1階層目だけでなく 2-3 階層目もそろえて、折りたたみの動作確認に使えるようにする。
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

    v_a1_c1_id INT;
    v_a1_c2_id INT;
    v_a1_c2_1_id INT;
    v_a2_c1_id INT;
    v_a2_c1_1_id INT;
    v_a3_c1_id INT;
    v_a3_c1_1_id INT;
    v_a3_c1_1_1_id INT;
    v_a4_c1_id INT;
    v_a5_c1_id INT;
    v_a5_c2_id INT;
BEGIN
    -- 冪等性: 本マイグレーション特有のコメント本文の存在で判定する
    IF EXISTS (
        SELECT 1 FROM replies
        WHERE content = '優勝チームのプレゼン、特に印象に残った工夫はありましたか？'
    ) THEN
        RETURN;
    END IF;

    SELECT id INTO v_article1_id FROM articles WHERE title = '神戸大学でのハッカソン体験記';
    SELECT id INTO v_article2_id FROM articles WHERE title = 'Goで作るREST API入門';
    SELECT id INTO v_article3_id FROM articles WHERE title = 'Vue 3 + Vuetifyで学ぶフロントエンド開発';
    SELECT id INTO v_article4_id FROM articles WHERE title = 'PostgreSQLのマイグレーション管理';
    SELECT id INTO v_article5_id FROM articles WHERE title = 'Dockerで開発環境を統一する';

    SELECT id INTO v_admin_id  FROM users WHERE name = 'admin';
    SELECT id INTO v_user01_id FROM users WHERE name = 'user01';
    SELECT id INTO v_user02_id FROM users WHERE name = 'user02';
    SELECT id INTO v_user03_id FROM users WHERE name = 'user03';

    -- ========== Article 1: 神戸大学でのハッカソン体験記 ==========
    -- 既存スレッドに加え、別系統のコメントを 2 本追加する
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_user03_id,
            '優勝チームのプレゼン、特に印象に残った工夫はありましたか？',
            'comment', NULL,
            '2026-04-04 09:30:00+09'::timestamptz,
            '2026-04-04 09:30:00+09'::timestamptz)
    RETURNING id INTO v_a1_c1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_admin_id,
            '実装より「誰の何の課題を解くか」を最初の3時間で詰めきったのが効いていました。',
            'comment', v_a1_c1_id,
            '2026-04-04 11:00:00+09'::timestamptz,
            '2026-04-04 11:00:00+09'::timestamptz);

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_user02_id,
            '次回参加するなら、事前に何を準備しておくと良さそうですか？',
            'comment', NULL,
            '2026-04-05 14:00:00+09'::timestamptz,
            '2026-04-05 14:00:00+09'::timestamptz)
    RETURNING id INTO v_a1_c2_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_admin_id,
            'チーム内でのGitフローと、よく使うUIコンポーネントの雛形だけ揃えておくと当日が楽です。',
            'comment', v_a1_c2_id,
            '2026-04-05 16:30:00+09'::timestamptz,
            '2026-04-05 16:30:00+09'::timestamptz)
    RETURNING id INTO v_a1_c2_1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article1_id, v_user02_id,
            'たしかに当日に Git で揉めるのはあるあるですね…参考にします！',
            'comment', v_a1_c2_1_id,
            '2026-04-05 18:00:00+09'::timestamptz,
            '2026-04-05 18:00:00+09'::timestamptz);

    -- ========== Article 2: Goで作るREST API入門 ==========
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article2_id, v_admin_id,
            'ミドルウェアの章、ロガーとリカバリの順序の説明が分かりやすかったです。',
            'comment', NULL,
            '2026-04-13 10:00:00+09'::timestamptz,
            '2026-04-13 10:00:00+09'::timestamptz)
    RETURNING id INTO v_a2_c1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article2_id, v_user01_id,
            'ありがとうございます！順序で挙動が変わるところは初学者がハマりやすいので意識して書きました。',
            'comment', v_a2_c1_id,
            '2026-04-13 12:00:00+09'::timestamptz,
            '2026-04-13 12:00:00+09'::timestamptz)
    RETURNING id INTO v_a2_c1_1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article2_id, v_user03_id,
            '横から失礼します。リカバリを外側にする派と内側にする派がいる気がしていて、考え方の違いを知りたいです。',
            'comment', v_a2_c1_1_id,
            '2026-04-13 13:30:00+09'::timestamptz,
            '2026-04-13 13:30:00+09'::timestamptz);

    -- ========== Article 3: Vue 3 + Vuetifyで学ぶフロントエンド開発 ==========
    -- 3 階層深いネストの動作確認用
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article3_id, v_user01_id,
            'props と emits をどこまで型で縛るべきか、いつも迷います。',
            'comment', NULL,
            '2026-04-17 09:00:00+09'::timestamptz,
            '2026-04-17 09:00:00+09'::timestamptz)
    RETURNING id INTO v_a3_c1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article3_id, v_user02_id,
            '原則 defineProps / defineEmits の型引数で書いておくと、エディタの補完が効いて事故が減ると思います。',
            'comment', v_a3_c1_id,
            '2026-04-17 10:30:00+09'::timestamptz,
            '2026-04-17 10:30:00+09'::timestamptz)
    RETURNING id INTO v_a3_c1_1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article3_id, v_user01_id,
            'ありがとうございます。defaults はどう書くのがおすすめですか？',
            'comment', v_a3_c1_1_id,
            '2026-04-17 11:00:00+09'::timestamptz,
            '2026-04-17 11:00:00+09'::timestamptz)
    RETURNING id INTO v_a3_c1_1_1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article3_id, v_user02_id,
            'withDefaults を併用すると素直に書けます。',
            'comment', v_a3_c1_1_1_id,
            '2026-04-17 11:20:00+09'::timestamptz,
            '2026-04-17 11:20:00+09'::timestamptz);

    -- ========== Article 4: PostgreSQLのマイグレーション管理 ==========
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article4_id, v_admin_id,
            'down マイグレーションをどこまで真面目に書くべきか、チームでも議論になります。',
            'comment', NULL,
            '2026-04-19 09:00:00+09'::timestamptz,
            '2026-04-19 09:00:00+09'::timestamptz)
    RETURNING id INTO v_a4_c1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article4_id, v_user03_id,
            '本番では基本 up しか流さない前提なので、down は「ローカルで巻き戻せる程度」で割り切る運用にしています。',
            'comment', v_a4_c1_id,
            '2026-04-19 10:00:00+09'::timestamptz,
            '2026-04-19 10:00:00+09'::timestamptz);

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article4_id, v_user02_id,
            '記事の例どおりに動かせました！seed もマイグレーションで管理する流派は新鮮でした。',
            'comment', NULL,
            '2026-04-19 13:30:00+09'::timestamptz,
            '2026-04-19 13:30:00+09'::timestamptz);

    -- ========== Article 5: Dockerで開発環境を統一する ==========
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article5_id, v_user02_id,
            'volumes の指定で、Mac だけ妙に遅くなる問題に最近ぶつかりました。',
            'comment', NULL,
            '2026-04-21 09:00:00+09'::timestamptz,
            '2026-04-21 09:00:00+09'::timestamptz)
    RETURNING id INTO v_a5_c1_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article5_id, v_user01_id,
            'node_modules を named volume に逃がすと体感かなり改善しますよ。',
            'comment', v_a5_c1_id,
            '2026-04-21 10:00:00+09'::timestamptz,
            '2026-04-21 10:00:00+09'::timestamptz);

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article5_id, v_admin_id,
            'compose ファイルに env_file を分けておくと、本番想定の検証も同じ手順で回せて便利でした。',
            'comment', NULL,
            '2026-04-21 14:00:00+09'::timestamptz,
            '2026-04-21 14:00:00+09'::timestamptz)
    RETURNING id INTO v_a5_c2_id;

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_article5_id, v_user03_id,
            'env_file 分割いいですね、真似します！',
            'comment', v_a5_c2_id,
            '2026-04-21 15:00:00+09'::timestamptz,
            '2026-04-21 15:00:00+09'::timestamptz);
END $$;
