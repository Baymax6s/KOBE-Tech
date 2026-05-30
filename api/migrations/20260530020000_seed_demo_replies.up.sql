-- 追加記事に対する返信（コメント / 質問 / 回答）の seed。
-- kind は VARCHAR で 'comment' / 'question' / 'answer'。
-- ネストは maxReplyDepth = 3（ルート depth0 〜 最大 depth3 の 4 階層）を超えないようにする。
--
-- 網羅する状態:
--   - 4 階層フルのコメントスレッド（折りたたみ確認用）            … 記事1
--   - ベストアンサーありの質問＋経路外の別回答                    … 記事1
--   - ベストアンサーありの質問（回答1件）                        … 記事2
--   - 回答はあるがベスト未選択の質問                             … 記事3
--   - 回答ゼロの未回答質問                                      … 記事4
--   - コメント無し（空スレッド）                                … 記事5 / 記事6（何も入れない）
DO $$
DECLARE
    v_art_git   INT;
    v_art_ts    INT;
    v_art_book  INT;
    v_art_css   INT;

    v_suzuki    INT; -- 鈴木健一
    v_takahashi INT; -- 高橋美咲
    v_watanabe  INT; -- 渡辺翔太
    v_ito       INT; -- 伊藤さくら
    v_nakamura  INT; -- 中村大輔
    v_tanaka    INT; -- 田中太郎
    v_yamada    INT; -- 山田花子

    v_c1   INT;
    v_c2   INT;
    v_c3   INT;
    v_q1   INT;
    v_q3   INT;
    v_book_c2 INT;
BEGIN
    -- 冪等性: 本マイグレーション固有のコメント本文の存在で判定する
    IF EXISTS (
        SELECT 1 FROM replies
        WHERE content = 'rebase、いつも怖くて避けていました…どこから慣れるのが良いですか？'
    ) THEN
        RETURN;
    END IF;

    SELECT id INTO v_art_git  FROM articles WHERE title = 'Gitのrebaseとmergeを使い分ける';
    SELECT id INTO v_art_ts   FROM articles WHERE title = 'TypeScriptの型で安全なAPIクライアントを作る';
    SELECT id INTO v_art_book FROM articles WHERE title = '個人開発した蔵書管理アプリを公開します';
    SELECT id INTO v_art_css  FROM articles WHERE title = 'CSS Grid と Flexbox の使い分けメモ';

    SELECT id INTO v_suzuki    FROM users WHERE name = '鈴木健一';
    SELECT id INTO v_takahashi FROM users WHERE name = '高橋美咲';
    SELECT id INTO v_watanabe  FROM users WHERE name = '渡辺翔太';
    SELECT id INTO v_ito       FROM users WHERE name = '伊藤さくら';
    SELECT id INTO v_nakamura  FROM users WHERE name = '中村大輔';
    SELECT id INTO v_tanaka    FROM users WHERE name = '田中太郎';
    SELECT id INTO v_yamada    FROM users WHERE name = '山田花子';

    -- ========== 記事1: Git — 4階層フルのコメントスレッド ==========
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_git, v_takahashi,
            'rebase、いつも怖くて避けていました…どこから慣れるのが良いですか？',
            'comment', NULL,
            '2026-05-06 12:00:00+09'::timestamptz, '2026-05-06 12:00:00+09'::timestamptz)
    RETURNING id INTO v_c1; -- depth 0

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_git, v_suzuki,
            '最初は merge だけでも大丈夫ですよ。慣れてきたら自分の作業ブランチで試すのがおすすめです。',
            'comment', v_c1,
            '2026-05-06 13:00:00+09'::timestamptz, '2026-05-06 13:00:00+09'::timestamptz)
    RETURNING id INTO v_c2; -- depth 1

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_git, v_takahashi,
            'rebase したあとに force push が必要になる場面がよく分かっていません。',
            'comment', v_c2,
            '2026-05-06 14:00:00+09'::timestamptz, '2026-05-06 14:00:00+09'::timestamptz)
    RETURNING id INTO v_c3; -- depth 2

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_git, v_suzuki,
            'rebase は履歴を書き換えるので、push 済みのブランチだと force が要ります。共有ブランチでは避けるのが安全です。',
            'comment', v_c3,
            '2026-05-06 15:00:00+09'::timestamptz, '2026-05-06 15:00:00+09'::timestamptz); -- depth 3（最深）

    -- ========== 記事1: Git — ベストアンサーあり + 経路外の別回答 ==========
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_git, v_nakamura,
            'コンフリクトが起きたとき、rebase の途中ではどう解消すればいいですか？',
            'question', NULL,
            '2026-05-07 09:00:00+09'::timestamptz, '2026-05-07 09:00:00+09'::timestamptz)
    RETURNING id INTO v_q1; -- depth 0

    INSERT INTO replies (article_id, user_id, content, is_best, kind, parent_id, created_at, updated_at)
    VALUES (v_art_git, v_suzuki,
            '衝突したファイルを直して git add し、git rebase --continue で先へ進めます。これを衝突がなくなるまで繰り返します。',
            TRUE, 'answer', v_q1,
            '2026-05-07 10:30:00+09'::timestamptz, '2026-05-07 10:30:00+09'::timestamptz); -- best, depth1

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_git, v_watanabe,
            '途中でやめたくなったら git rebase --abort で元の状態に戻せますよ。',
            'answer', v_q1,
            '2026-05-07 11:00:00+09'::timestamptz, '2026-05-07 11:00:00+09'::timestamptz); -- 経路外の兄弟, depth1

    -- ========== 記事2: TypeScript — コメント + ベストアンサーあり ==========
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_ts, v_watanabe,
            'ジェネリクスでラップする書き方、ちょうど知りたかったので助かりました！',
            'comment', NULL,
            '2026-05-11 09:00:00+09'::timestamptz, '2026-05-11 09:00:00+09'::timestamptz)
    RETURNING id INTO v_c1; -- depth 0

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_ts, v_takahashi,
            'お役に立ててうれしいです！戻り値の型を1か所に集約できるのがポイントです。',
            'comment', v_c1,
            '2026-05-11 10:00:00+09'::timestamptz, '2026-05-11 10:00:00+09'::timestamptz); -- depth 1

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_ts, v_ito,
            '返ってくる JSON のうち一部の項目だけ型を付けたいときは、どう書くのが良いでしょう？',
            'question', NULL,
            '2026-05-12 13:00:00+09'::timestamptz, '2026-05-12 13:00:00+09'::timestamptz)
    RETURNING id INTO v_q1; -- depth 0

    INSERT INTO replies (article_id, user_id, content, is_best, kind, parent_id, created_at, updated_at)
    VALUES (v_art_ts, v_takahashi,
            '使う項目だけを並べた型を定義して、それで受け取れば十分です。余分な項目は無視されます。',
            TRUE, 'answer', v_q1,
            '2026-05-12 15:00:00+09'::timestamptz, '2026-05-12 15:00:00+09'::timestamptz); -- best, depth1

    -- ========== 記事3: 蔵書管理アプリ — 称賛コメント + ベスト未選択のQ&A ==========
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_book, v_takahashi,
            '読了数をグラフで見られるの、モチベが上がりそうで良いですね！',
            'comment', NULL,
            '2026-05-15 16:00:00+09'::timestamptz, '2026-05-15 16:00:00+09'::timestamptz);

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_book, v_suzuki,
            '授業の構成をそのまま個人開発に持ち込むの、とても良い流れですね。',
            'comment', NULL,
            '2026-05-16 09:00:00+09'::timestamptz, '2026-05-16 09:00:00+09'::timestamptz)
    RETURNING id INTO v_book_c2; -- depth 0

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_book, v_watanabe,
            'ありがとうございます！感想メモ機能も近いうちに足す予定です。',
            'comment', v_book_c2,
            '2026-05-16 10:00:00+09'::timestamptz, '2026-05-16 10:00:00+09'::timestamptz); -- depth 1

    -- ベスト未選択の質問（回答は2件あるが、どれも is_best = FALSE）
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_book, v_yamada,
            '書影の取得には、どの API を使っていますか？',
            'question', NULL,
            '2026-05-17 11:00:00+09'::timestamptz, '2026-05-17 11:00:00+09'::timestamptz)
    RETURNING id INTO v_q3; -- depth 0

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_book, v_nakamura,
            'openBD が無料で使いやすいと聞いたことがあります。',
            'answer', v_q3,
            '2026-05-17 12:00:00+09'::timestamptz, '2026-05-17 12:00:00+09'::timestamptz); -- depth1, not best

    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_book, v_tanaka,
            'Google Books API も候補になりますよ。表紙画像の有無は本によりますが。',
            'answer', v_q3,
            '2026-05-17 13:30:00+09'::timestamptz, '2026-05-17 13:30:00+09'::timestamptz); -- depth1, not best

    -- ========== 記事4: CSS — 少コメント + 未回答の質問 ==========
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_css, v_takahashi,
            '私もいつも迷うので、この方針メモ助かります！',
            'comment', NULL,
            '2026-05-19 18:00:00+09'::timestamptz, '2026-05-19 18:00:00+09'::timestamptz);

    -- 未回答の質問（子の回答なし）
    INSERT INTO replies (article_id, user_id, content, kind, parent_id, created_at, updated_at)
    VALUES (v_art_css, v_watanabe,
            'grid-template-columns の auto-fit と auto-fill って、何が違うんでしょうか？',
            'question', NULL,
            '2026-05-20 10:00:00+09'::timestamptz, '2026-05-20 10:00:00+09'::timestamptz);

    -- 記事5（学習習慣）と記事6（最新）はコメント無しの状態を見せるため、ここでは何も追加しない。
END $$;
