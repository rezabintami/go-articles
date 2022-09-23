CREATE TABLE IF NOT EXISTS comments (
        id bigserial PRIMARY KEY,
        user_id integer NOT NULL,
        comment TEXT NOT NULL,
        created_at timestamptz NOT NULL,
        updated_at timestamptz NULL,
        deleted_at timestamptz NULL
    );