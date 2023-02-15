-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  first_name VARCHAR(256) NOT NULL,
  last_name VARCHAR(256) NOT NULL,
  nickname VARCHAR(256) NOT NULL,
  password VARCHAR(256) NOT NULL,
  email VARCHAR(256) NOT NULL UNIQUE,
  country VARCHAR(256) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_country ON users (country);

-- +goose Down
DROP TABLE users;
DROP INDEX IF EXISTS idx_users_country;
