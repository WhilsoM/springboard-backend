package dto

import "time"

type OpportunityResponse struct {
	ID              string     `json:"id"`
	EmployerID      string     `json:"employer_id"`
	CompanyName     string     `json:"company_name"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Type            string     `json:"type"`
	Format          string     `json:"format"`
	City            string     `json:"city"`
	Address         string     `json:"address"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
	Tags            []string   `json:"tags"`
	SalaryMin       *int       `json:"salary_min,omitempty"`
	SalaryMax       *int       `json:"salary_max,omitempty"`
	ExperienceLevel string     `json:"experience_level,omitempty"`
	IsActive        bool       `json:"is_active"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type CreateOpportunityRequest struct {
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Type            string     `json:"type"`
	Format          string     `json:"format"`
	City            string     `json:"city"`
	Address         string     `json:"address"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
	Tags            []string   `json:"tags"`
	SalaryMin       *int       `json:"salary_min"`
	SalaryMax       *int       `json:"salary_max"`
	ExperienceLevel string     `json:"experience_level"`
	ExpiresAt       *time.Time `json:"expires_at"`
}

type SearchFilters struct {
	Type   string   `json:"type"`
	Format string   `json:"format"`
	City   string   `json:"city"`
	Tags   []string `json:"tags"`
	Search string   `json:"search"`
}
