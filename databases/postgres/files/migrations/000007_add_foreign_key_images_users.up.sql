ALTER TABLE 
    users 
ADD
    CONSTRAINT fk_images_users FOREIGN KEY(image_id) REFERENCES images(id);
