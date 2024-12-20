CREATE TABLE files_actions_log (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON UPDATE CASCADE,
    file_id bigint NOT NULL REFERENCES files(id) ON UPDATE CASCADE,
    action_id bigint NOT NULL REFERENCES files_actions(id) ON UPDATE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);