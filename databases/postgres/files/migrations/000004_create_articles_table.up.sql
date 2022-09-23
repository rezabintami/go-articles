CREATE TABLE IF NOT EXISTS articles (
    id bigserial PRIMARY KEY,
    user_id integer NOT NULL,
    image_id integer NULL UNIQUE, 
    title varchar(255) NOT NULL,
    description TEXT NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL 
);