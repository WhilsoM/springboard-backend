-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION create_profile_after_user_insert()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.role = 'applicant' THEN
        INSERT INTO candidate_profiles (user_id)
        VALUES (NEW.id);
    ELSIF NEW.role = 'employer' THEN
        INSERT INTO employer_profiles (user_id, company_name)
        VALUES (NEW.id, NEW.display_name);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trigger_create_profile
AFTER INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION create_profile_after_user_insert();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_create_profile ON users;
DROP FUNCTION IF EXISTS create_profile_after_user_insert();
