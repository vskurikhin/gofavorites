-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_created_at_favorites ON favorites;
CREATE TRIGGER set_created_at_favorites
    BEFORE INSERT
    ON favorites
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_favorites ON favorites;
CREATE TRIGGER set_update_at_favorites
    BEFORE UPDATE
    ON favorites
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_update_at_favorites ON users;
DROP TRIGGER IF EXISTS set_created_at_favorites ON users;
-- +goose StatementEnd
