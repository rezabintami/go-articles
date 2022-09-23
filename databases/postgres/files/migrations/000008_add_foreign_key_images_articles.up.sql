ALTER TABLE 
    articles
ADD
    CONSTRAINT fk_images_articles FOREIGN KEY(image_id) REFERENCES images(id);