-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_types
(
    name       varchar PRIMARY KEY,
    deleted    bool,
    created_at timestamp NOT NULL,
    updated_at timestamp
);

COMMENT ON COLUMN asset_types.name IS 'Type of asset, example bonds, stocks, funds ... etc';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS asset_types;
-- +goose StatementEnd
