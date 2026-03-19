-- +goose Up
CREATE TABLE IF NOT EXISTS candidate_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    university TEXT,
    course INT CHECK (course > 0 AND course < 7),
    skills TEXT[] DEFAULT '{}',
    portfolio_url TEXT,
    github_url TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_candidate_skills ON candidate_profiles USING GIN (skills);

-- +goose Down
DROP TABLE IF EXISTS candidate_profiles;
