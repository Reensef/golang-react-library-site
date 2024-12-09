CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS files_tags(
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL UNIQUE
);