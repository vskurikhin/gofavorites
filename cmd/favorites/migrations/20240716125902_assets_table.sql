-- +goose Up
-- +goose StatementBegin

CREATE TABLE assets
(
    isin       varchar PRIMARY KEY,
    asset_type varchar   NOT NULL,
    deleted    bool,
    created_at timestamp NOT NULL,
    updated_at timestamp
);

ALTER TABLE assets
    ADD FOREIGN KEY (asset_type)
        REFERENCES asset_types (name);

COMMENT ON COLUMN assets.isin IS 'International Securities Identification Numbers';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS assets;
-- +goose StatementEnd
