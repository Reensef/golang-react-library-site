CREATE TABLE files_actions (
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL
);

INSERT INTO files_actions (name) VALUES ('Downloaded');
INSERT INTO files_actions (name) VALUES ('Uploaded');
INSERT INTO files_actions (name) VALUES ('Opened');
INSERT INTO files_actions (name) VALUES ('Deleted');