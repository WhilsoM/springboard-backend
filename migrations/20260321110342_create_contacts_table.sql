-- +goose Up
CREATE TYPE contact_status AS ENUM ('pending', 'accepted', 'rejected');

CREATE TABLE IF NOT EXISTS contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status contact_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(sender_id, receiver_id)
);

CREATE INDEX idx_contacts_users ON contacts(sender_id, receiver_id);

-- +goose Down
DROP TABLE IF EXISTS contacts;
DROP TYPE IF EXISTS contact_status;
