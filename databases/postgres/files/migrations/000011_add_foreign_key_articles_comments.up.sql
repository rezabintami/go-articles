ALTER TABLE articles_comments
ADD
    CONSTRAINT fk_articles_articles_comments FOREIGN KEY(article_id) REFERENCES articles(id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE articles_comments 
ADD
    CONSTRAINT fk_comments_articles_comments FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE ON UPDATE CASCADE;