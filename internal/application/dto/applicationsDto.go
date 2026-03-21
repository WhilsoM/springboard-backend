package dto

import "time"

type ApplyRequest struct {
	CoverLetter string `json:"cover_letter"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"` // pending, accepted, rejected, reserve
}

type PaginationFilters struct {
	Limit  int
	Offset int
}

// data from candidate_profiles and users tables
type ApplicantInfo struct {
	ID          string   `json:"id"`
	DisplayName string   `json:"display_name"`
	University  string   `json:"university"`
	Course      int      `json:"course"`
	Skills      []string `json:"skills"`
	GitHubURL   string   `json:"github_url"`
	AvatarURL   string   `json:"avatar_url"`
}

type ApplicationResponse struct {
	ID            string         `json:"id"`
	OpportunityID string         `json:"opportunity_id"`
	Status        string         `json:"status"`
	CoverLetter   string         `json:"cover_letter,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	Applicant     *ApplicantInfo `json:"applicant,omitempty"`
}
