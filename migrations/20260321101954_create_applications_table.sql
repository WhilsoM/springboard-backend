-- +goose Up
CREATE TYPE application_status AS ENUM ('pending', 'accepted', 'rejected', 'reserve');

CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    opportunity_id UUID NOT NULL REFERENCES opportunities(id) ON DELETE CASCADE,
    applicant_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cover_letter TEXT,
    status application_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(opportunity_id, applicant_id)
);

CREATE INDEX idx_apps_opportunity ON applications(opportunity_id);
CREATE INDEX idx_apps_applicant ON applications(applicant_id);

-- +goose Down
DROP TABLE IF EXISTS applications;
DROP TYPE IF EXISTS application_status;
