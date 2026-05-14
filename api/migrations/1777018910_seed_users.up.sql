
INSERT INTO users (id, name, password_hash) VALUES
(1, 'admin',  '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
(2, 'user01', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
(3, 'user02', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
(4, 'user03', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u')
ON CONFLICT (name) DO NOTHING;
