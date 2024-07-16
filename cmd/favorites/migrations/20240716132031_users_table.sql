-- +goose Up
-- +goose StatementBegin

CREATE TABLE users
(
    upk        varchar PRIMARY KEY,
    deleted    bool,
    created_at timestamp NOT NULL,
    updated_at timestamp
);

COMMENT ON COLUMN users.upk IS 'User primary key';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
