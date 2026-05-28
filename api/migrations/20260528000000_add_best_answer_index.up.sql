CREATE INDEX idx_replies_user_best ON replies(user_id, article_id) WHERE is_best = TRUE;
