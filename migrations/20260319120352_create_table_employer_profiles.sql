-- +goose Up
CREATE TABLE IF NOT EXISTS employer_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    company_name TEXT NOT NULL,
    inn VARCHAR(12) UNIQUE,
    description TEXT,
    website_url TEXT,
    logo_url TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_employer_company_name ON employer_profiles (company_name);

-- +goose Down
DROP TABLE IF EXISTS employer_profiles;
