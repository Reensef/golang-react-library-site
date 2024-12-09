CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS files(
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    size bigint NOT NULL,
    created_by_user_id bigserial NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user
        FOREIGN KEY (created_by_user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);