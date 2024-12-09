CREATE TABLE IF NOT EXISTS file_to_tags (
    file_id bigserial NOT NULL,
    tag_id bigserial NOT NULL,
    PRIMARY KEY (file_id, tag_id),
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES files_tags(id) ON DELETE CASCADE
);