CREATE TABLE IF NOT EXISTS posts(
                                id uuid PRIMARY KEY,
                                author_id uuid,
                                post_text TEXT,
                                published TIMESTAMP,
                                disable_comments BOOLEAN,
                                CONSTRAINT fk_author
                                    FOREIGN KEY(author_id)
                                        REFERENCES users(id)
);