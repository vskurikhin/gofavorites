-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_update_at() RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$$
BEGIN
    NEW.created_at = OLD.created_at;
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_update_at();
-- +goose StatementEnd
