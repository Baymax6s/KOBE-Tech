CREATE TABLE replies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    article_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    is_best BOOLEAN NOT NULL DEFAULT FALSE,
    parent_id INT NULL,
    kind TINYINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_replies_article
        FOREIGN KEY (article_id) REFERENCES articles(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_replies_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_replies_parent
        FOREIGN KEY (parent_id) REFERENCES replies(id)
        ON DELETE CASCADE,
    CONSTRAINT chk_replies_kind
        CHECK (kind IN (0, 1, 2))
);

CREATE INDEX idx_replies_article_id ON replies(article_id);
CREATE INDEX idx_replies_user_id ON replies(user_id);
CREATE INDEX idx_replies_parent_id ON replies(parent_id);