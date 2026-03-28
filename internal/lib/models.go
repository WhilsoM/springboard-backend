package lib

import (
	"time"
)

type User struct {
	PasswordHash string   `json:"-" db:"password_hash"`
	ID           string   `json:"id" db:"id"`
	Role         UserRole `json:"role" db:"role"`
	Email        string   `json:"email" db:"email"`
	DisplayName  string   `json:"display_name" db:"display_name"`
}

type ApplicantUser struct {
	User
	University   *string   `json:"university"`
	Course       int       `json:"course"`
	Skills       []string  `json:"skills"`
	PortfolioURL *string   `json:"portfolio_url"`
	GithubURL    *string   `json:"github_url"`
	AvatarURL    *string   `json:"avatar_url"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EmployerUser struct {
	User
	CompanyName string    `json:"company_name"`
	IsVerified  bool      `json:"is_verified"`
	INN         *string   `json:"inn"`
	Description *string   `json:"description"`
	AvatarURL   *string   `json:"avatar_url"`
	WebsiteURL  *string   `json:"website_url"`
	UpdatedAt   time.Time `json:"updated_at"`
}
