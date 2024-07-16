-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_created_at() RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$$
BEGIN
    NEW.created_at = now();
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_created_at();
-- +goose StatementEnd
