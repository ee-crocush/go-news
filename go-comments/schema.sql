DROP TABLE IF EXISTS comments;
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    parent_id BIGINT,
    user_name TEXT NOT NULL,
    content TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0
);