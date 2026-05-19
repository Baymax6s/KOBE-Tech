-- Article 1 のベストアンサー経路の外側にも回答を 2 件追加する。
-- 「ベストアンサー経路上のノードだけが初期表示される」UI を確認するには
-- 経路外のノード（=初期非表示になるべきもの）が seed に存在している必要があるため。
--   - ルート直下の兄弟（経路から最も離れた位置）
--   - 経路途中の追問の兄弟（ベスト直近の位置）
-- の 2 箇所に配置することで、「N 件を表示」ボタンが現れること・件数合算が正しいことを目視できる。
DO $$
DECLARE
    v_article1_id  INT;
    v_user03_id    INT;
    v_question_id  INT; -- Article 1 のルート質問
    v_followup_id  INT; -- ルート質問配下の追問（ベストの親）
BEGIN
    -- 冪等性: 本マイグレーション特有の回答本文の存在で判定する
    IF EXISTS (
        SELECT 1 FROM replies
        WHERE content = '私はチームビルディング系の準備にも時間を割きました。役割分担と意思決定の練習をしておくと、技術面以上に効きます。'
    ) THEN
        RETURN;
    END IF;

    SELECT id INTO v_article1_id FROM articles WHERE title = '神戸大学でのハッカソン体験記';
    SELECT id INTO v_user03_id   FROM users    WHERE name  = 'user03';

    SELECT id INTO v_question_id FROM replies
    WHERE article_id = v_article1_id
      AND content    = '次回参加するときに、当日までに練習しておくと良い技術スタックはありますか？';

    SELECT id INTO v_followup_id FROM replies
    WHERE article_id = v_article1_id
      AND content    = 'なるほど、フロントとバックの最小サンプルというのは、具体的にはどんな構成でしょうか？';

    -- 先行 seed が未適用な環境では親が見つからない可能性があるのでスキップする。
    IF v_question_id IS NULL OR v_followup_id IS NULL THEN
        RAISE NOTICE 'seed_off_path_replies: parent reply not found, skipping.';
        RETURN;
    END IF;

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
