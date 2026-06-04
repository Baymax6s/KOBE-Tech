-- 本マイグレーションで追加した記事・タグ紐付け・追加タグを削除する。
-- 返信・いいねは後続マイグレーション (020000/030000) の down で先に消える前提。

-- 記事 × タグの紐付けを外す
DELETE FROM article_tags at
USING articles a
WHERE at.article_id = a.id
  AND a.title IN (
        'Gitのrebaseとmergeを使い分ける',
        'TypeScriptの型で安全なAPIクライアントを作る',
        '個人開発した蔵書管理アプリを公開します',
        'CSS Grid と Flexbox の使い分けメモ',
        'プログラミング学習を続けるための小さな習慣',
        'Postmanで始めるAPIデバッグ入門',
        'C言語課題36：ポインタのイメージを掴む',
        'C#課題24：例外処理でよくある詰まり',
        'Java課題11：継承とインターフェースの使い分け'
    );

-- 記事本体を削除
DELETE FROM articles
WHERE title IN (
    'Gitのrebaseとmergeを使い分ける',
    'TypeScriptの型で安全なAPIクライアントを作る',
    '個人開発した蔵書管理アプリを公開します',
    'CSS Grid と Flexbox の使い分けメモ',
    'プログラミング学習を続けるための小さな習慣',
    'Postmanで始めるAPIデバッグ入門',
    'C言語課題36：ポインタのイメージを掴む',
    'C#課題24：例外処理でよくある詰まり',
    'Java課題11：継承とインターフェースの使い分け'
);

-- 本マイグレーションで追加したタグを削除（他記事で使われていないもののみ安全に消す）
DELETE FROM tags t
WHERE t.name IN (
        'Git', 'TypeScript', 'CSS', '個人開発', '制作物',
        'C言語課題36', 'C#課題24', 'Java課題11', 'プロジェクト管理', 'デザイン思考'
    )
  AND NOT EXISTS (
        SELECT 1 FROM article_tags at WHERE at.tag_id = t.id
  );
