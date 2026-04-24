INSERT INTO users (name, password_hash) VALUES
    ('admin',  '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
    ('user01', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
    ('user02', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u'),
    ('user03', '$2a$10$AiAy4O5udI.h/SeAzvtLF.mgpc07e9Xgb0V5teBs66Oxsnsjt603u')
ON CONFLICT (name) DO NOTHING;
