-- tags.name を大文字小文字を区別せずユニーク扱いするための部分式インデックス。
-- アプリ層では NormalizeTagNames が原文ケースを残すようになったので、"Vue" と "vue" を
-- 同一タグとして扱う制約は DB 側で担保する。
CREATE UNIQUE INDEX uniq_tags_lower_name ON tags (LOWER(name));
