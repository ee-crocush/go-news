CREATE TYPE comment_status AS ENUM ('pending', 'approved', 'rejected');
DROP TABLE IF EXISTS comments;
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    news_id INT NOT NULL,
    parent_id BIGINT,
    user_name TEXT NOT NULL,
    content TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0,
    status comment_status NOT NULL DEFAULT 'pending'
);