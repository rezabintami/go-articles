CREATE TABLE IF NOT EXISTS images (
        id bigserial PRIMARY KEY,
        path TEXT NOT NULL,
        type VARCHAR NOT NULL,
        created_at timestamptz NOT NULL
    );