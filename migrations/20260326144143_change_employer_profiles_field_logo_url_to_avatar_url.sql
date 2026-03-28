-- +goose Up
ALTER TABLE employer_profiles RENAME COLUMN logo_url TO avatar_url;

-- +goose Down
ALTER TABLE employer_profiles RENAME COLUMN avatar_url TO logo_url;
