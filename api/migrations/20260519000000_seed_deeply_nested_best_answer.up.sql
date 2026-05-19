-- ベストアンサーがネストの奥（質問の孫以降）にぶら下がるケースを 1 件だけ追加する。
-- 既存 seed では「質問の直下の回答がベスト」しかないため、
-- 「ベストアンサーへの経路上にある中間ノードも初期表示される」UI 挙動が目視確認できなかった。
-- Article 1 のスレッド末尾（追問で終わっている）の続きとして admin の解答を投入し、それをベストにする。
DO $$
DECLARE
    v_article1_id INT;
    v_admin_id    INT;
    v_parent_id   INT;
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

    -- 親となる「追問」reply を本文で特定する。
    -- 20260513 の seed が未適用な環境では親が見つからない可能性があるので、その場合はスキップする。
    SELECT id INTO v_parent_id FROM replies
    WHERE article_id = v_article1_id
      AND content    = 'なるほど、フロントとバックの最小サンプルというのは、具体的にはどんな構成でしょうか？';

    IF v_parent_id IS NULL THEN
        RAISE NOTICE 'seed_deeply_nested_best_answer: parent reply not found, skipping.';
        RETURN;
    END IF;

    INSERT INTO replies (article_id, user_id, content, is_best, kind, parent_id, created_at, updated_at)
    VALUES (
        v_article1_id, v_admin_id,
        '最小サンプルは「フロントが Hello を表示」「バックが GET /ping で pong を返す」「両者を fetch で繋ぐ」の 3 点セットで十分です。当日の認識合わせがぐっと早くなります。',
        TRUE, 2, v_parent_id,
        '2026-04-05 16:00:00+09'::timestamptz,
        '2026-04-05 16:00:00+09'::timestamptz
    );
END $$;
