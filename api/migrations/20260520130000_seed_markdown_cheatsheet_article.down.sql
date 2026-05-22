DELETE FROM articles
WHERE (title, user_id) IN (
    SELECT v.title, u.id
    FROM (VALUES
        ('Markdown 記法チートシート', 'admin')
    ) AS v(title, user_name)
    JOIN users u ON u.name = v.user_name
);
