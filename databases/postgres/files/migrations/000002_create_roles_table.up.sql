CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name varchar(255) UNIQUE NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NULL
);

INSERT INTO roles(id, name, created_at) VALUES(1, 'USER', NOW());
INSERT INTO roles(id, name, created_at) VALUES(2, 'ADMIN', NOW());
INSERT INTO roles(id, name, created_at) VALUES(3, 'SUPERUSER', NOW());