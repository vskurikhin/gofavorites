-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_created_at_users ON users;
CREATE TRIGGER set_created_at_users
    BEFORE INSERT
    ON users
    FOR EACH ROW
EXECUTE FUNCTION set_created_at();

DROP TRIGGER IF EXISTS set_update_at_users ON users;
CREATE TRIGGER set_update_at_users
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION set_update_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_update_at_users ON users;
DROP TRIGGER IF EXISTS set_created_at_users ON users;
-- +goose StatementEnd
