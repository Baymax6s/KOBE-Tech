-- Article 1 のスレッドに「ベストアンサーがネストの奥にぶら下がる」+「経路外にも回答がある」構造を作る。
-- 「ベストアンサー経路上のノードだけが初期表示される」UI を目視確認するには、
-- 経路上の中間ノード（孫以降にぶら下がるベスト）と、経路外のノード（隠れるべき兄弟）の
-- 両方が seed に必要なため。
--
-- 追加後のスレッド構造:
--   質問(山田花子)
--   ├─ 回答(admin)             ← 経路上（既存 seed）
--   │  └─ 追問(山田花子)         ← 経路上（既存 seed）
--   │     ├─ 回答(佐藤次郎) NEW  ← 経路外: ベスト直近の兄弟
--   │     └─ 回答(admin) NEW   ← 経路の終端（ベストアンサー）
--   └─ 回答(佐藤次郎) NEW         ← 経路外: ルート直下の兄弟
DO $$
DECLARE
    v_article1_id INT;
    v_admin_id    INT;
    v_user03_id   INT;
    v_question_id INT; -- ルート質問の id
    v_followup_id INT; -- 追問（=ベストの親）の id
BEGIN
    -- 冪等性: 本マイグレーション特有の回答本文の存在で判定する
    IF EXISTS (
        SELECT 1 FROM replies
        WHERE content = '最小サンプルは「フロントが Hello を表示」「バックが GET /ping で pong を返す」「両者を fetch で繋ぐ」の 3 点セットで十分です。当日の認識合わせがぐっと早くなります。'
    ) THEN
        RETURN;
    END IF;

    SELECT id INTO v_article1_id FROM articles WHERE title = '神戸大学でのハッカソン体験記';
    SELECT id INTO v_admin_id    FROM users    WHERE name  = 'admin';
    SELECT id INTO v_user03_id   FROM users    WHERE name  = '佐藤次郎';

    SELECT id INTO v_question_id FROM replies
    WHERE article_id = v_article1_id
      AND content    = '次回参加するときに、当日までに練習しておくと良い技術スタックはありますか？';

    SELECT id INTO v_followup_id FROM replies
    WHERE article_id = v_article1_id
      AND content    = 'なるほど、フロントとバックの最小サンプルというのは、具体的にはどんな構成でしょうか？';

    -- 先行 seed (20260513) が未適用な環境では親が見つからない可能性があるのでスキップする。
    IF v_question_id IS NULL OR v_followup_id IS NULL THEN
        RAISE NOTICE 'seed_deeply_nested_best_answer: parent reply not found, skipping.';
        RETURN;
    END IF;

    -- 20260518 のマイグレーションで kind は VARCHAR ('comment' / 'question' / 'answer') に変更済み。
    -- 整数リテラル (1,2) では chk_replies_kind に弾かれるので、必ず文字列で渡す。

    -- 経路の終端: 追問の下にベストアンサーを追加する
    INSERT INTO replies (article_id, user_id, content, is_best, kind, parent_id, created_at, updated_at)
    VALUES (
        v_article1_id, v_admin_id,
        '最小サンプルは「フロントが Hello を表示」「バックが GET /ping で pong を返す」「両者を fetch で繋ぐ」の 3 点セットで十分です。当日の認識合わせがぐっと早くなります。',
        TRUE, 'answer', v_followup_id,
        '2026-04-05 16:00:00+09'::timestamptz,
        '2026-04-05 16:00:00+09'::timestamptz
    );

    -- 経路外 (1): ルート質問の別兄弟回答（経路から最も離れた位置）
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (
        v_article1_id, v_user03_id,
        '私はチームビルディング系の準備にも時間を割きました。役割分担と意思決定の練習をしておくと、技術面以上に効きます。',
        'answer', v_question_id,
        '2026-04-05 14:00:00+09'::timestamptz,
        '2026-04-05 14:00:00+09'::timestamptz
    );

    -- 経路外 (2): 追問の別兄弟回答（ベストアンサーのすぐ隣）
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (
        v_article1_id, v_user03_id,
        '補足ですが、DB の最小サンプル（例: PostgreSQL に 1 テーブル）も用意しておくと、当日 API と DB の繋ぎ込み練習がそのまま流用できます。',
        'answer', v_followup_id,
        '2026-04-05 15:00:00+09'::timestamptz,
        '2026-04-05 15:00:00+09'::timestamptz
    );
END $$;
