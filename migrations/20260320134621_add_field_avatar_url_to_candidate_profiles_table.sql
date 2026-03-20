-- +goose Up
ALTER TABLE candidate_profiles ADD COLUMN avatar_url TEXT;

-- +goose Down
ALTER TABLE candidate_profiles DROP COLUMN IF EXISTS avatar_url;
