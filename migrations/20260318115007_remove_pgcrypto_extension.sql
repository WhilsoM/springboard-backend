-- +goose Up
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "pgcrypto";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- +goose StatementEnd
