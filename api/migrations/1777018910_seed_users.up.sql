
INSERT INTO users (id, name, password_hash) VALUES
(1, 'admin',  '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
(2, '田中太郎', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
(3, '山田花子', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
(4, '佐藤次郎', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u')
ON CONFLICT (name) DO NOTHING;
