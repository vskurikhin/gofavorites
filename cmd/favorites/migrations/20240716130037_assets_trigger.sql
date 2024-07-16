-- +goose Up
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_created_at_assets ON assets;
CREATE TRIGGER set_created_at_assets
    BEFORE INSERT
    ON assets
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_assets ON assets;
CREATE TRIGGER set_update_at_assets
    BEFORE UPDATE
    ON assets
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_update_at_assets ON assets;
DROP TRIGGER IF EXISTS set_created_at_assets ON assets;
-- +goose StatementEnd
