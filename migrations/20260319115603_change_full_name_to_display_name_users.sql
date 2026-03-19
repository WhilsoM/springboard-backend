-- +goose Up
ALTER TABLE users RENAME COLUMN full_name TO display_name;

-- +goose Down
ALTER TABLE users RENAME COLUMN display_name TO full_name;
