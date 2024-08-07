-- +goose Up
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_created_at_asset_types ON asset_types;
CREATE TRIGGER set_created_at_asset_types
    BEFORE INSERT
    ON asset_types
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_asset_types ON asset_types;
CREATE TRIGGER set_update_at_asset_types
    BEFORE UPDATE
    ON asset_types
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_update_at_asset_types ON asset_types;
DROP TRIGGER IF EXISTS set_created_at_asset_types ON asset_types;
-- +goose StatementEnd
