-- +goose Up
-- +goose StatementBegin

CREATE TABLE favorites
(
    id         uuid PRIMARY KEY DEFAULT pg_catalog.uuid_generate_v4(),
    isin       varchar   NOT NULL,
    user_upk   varchar   NOT NULL,
    version    bigint,
    deleted    bool,
    created_at timestamp NOT NULL,
    updated_at timestamp
);

CREATE UNIQUE INDEX IF NOT EXISTS favorites_bkey
    ON favorites (isin, user_upk);

ALTER TABLE favorites
    ADD FOREIGN KEY (user_upk)
        REFERENCES users (upk);

ALTER TABLE favorites
    ADD FOREIGN KEY (isin)
        REFERENCES assets (isin);

COMMENT ON COLUMN favorites.version IS 'CBUL version';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS favorites;
-- +goose StatementEnd
