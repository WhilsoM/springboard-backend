-- +goose Up
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO tags (name) VALUES
('Java'), ('Python'), ('Go'), ('React'), ('Rust'), ('Vue'), ('Junior'), ('Middle'), ('Part-time'), ('Remote')
ON CONFLICT DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS tags;
