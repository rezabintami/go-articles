CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  role_id integer NOT NULL, 
  image_id integer NULL UNIQUE, 
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NULL
);