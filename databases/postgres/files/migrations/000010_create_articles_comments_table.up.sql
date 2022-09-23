CREATE TABLE IF NOT EXISTS articles_comments (
    id bigserial PRIMARY KEY,
    article_id integer NOT NULL, 
    comment_id integer NOT NULL 
)