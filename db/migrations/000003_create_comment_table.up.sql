CREATE TABLE IF NOT EXISTS comments
(
    id        uuid PRIMARY KEY,
    author_id uuid,
    published TIMESTAMP,
    text varchar(2000),
    post_id uuid,
    parent_id uuid,
    CONSTRAINT fk_author
        FOREIGN KEY(author_id)
            REFERENCES users(id),
    CONSTRAINT fk_post
        FOREIGN KEY(post_id)
            REFERENCES posts(id),
    CONSTRAINT fk_parent
        FOREIGN KEY(parent_id)
            REFERENCES comments(id)
)