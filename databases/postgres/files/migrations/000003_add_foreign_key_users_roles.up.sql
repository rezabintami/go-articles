ALTER TABLE
    users
ADD
    CONSTRAINT fk_users_roles FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE;
