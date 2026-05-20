-- tags.name を大文字小文字を区別せずユニーク扱いするための部分式インデックス。
-- アプリ層では NormalizeTagNames が原文ケースを残すようになったので、"Vue" と "vue" を
-- 同一タグとして扱う制約は DB 側で担保する。
CREATE UNIQUE INDEX uniq_tags_lower_name ON tags (LOWER(name));

-- 旧 UNIQUE (name) は uniq_tags_lower_name に含意される（同じ name は同じ LOWER に潰れるため）。
-- 制約が 2 つ並ぶとどちらが効いているのか分かりにくく、INSERT のインデックス更新も二重になるので削除。
ALTER TABLE tags DROP CONSTRAINT tags_name_key;
