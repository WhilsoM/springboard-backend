-- +goose Up
CREATE TABLE IF NOT EXISTS opportunities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employer_id UUID NOT NULL REFERENCES employer_profiles(user_id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('internship', 'vacancy', 'event', 'mentorship')),
    format TEXT NOT NULL CHECK (format IN ('office', 'hybrid', 'remote')),
    city TEXT NOT NULL,
    address TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    tags TEXT[] DEFAULT '{}',
    salary_min INT,
    salary_max INT,
    experience_level TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_opp_type ON opportunities(type);
CREATE INDEX idx_opp_tags ON opportunities USING GIN (tags);
CREATE INDEX idx_opp_location ON opportunities(city);

-- +goose Down
DROP TABLE IF EXISTS opportunities;
