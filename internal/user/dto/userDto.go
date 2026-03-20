package dto

type UserResponse struct {
	Data any `json:"user"`
}

type updateMe struct {
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	DisplayName string `json:"display_name" db:"display_name"`
	IsPrivate   bool   `json:"is_private" db:"is_private"`
}

type UpdateMeEmployerRequest struct {
	updateMe
	CompanyName string `json:"company_name" db:"company_name"`
	Description string `json:"description" db:"description"`
	WebsiteURL  string `json:"website_url" db:"website_url"`
	LogoURL     string `json:"logo_url" db:"logo_url"`
	IsVerified  bool   `json:"is_verified" db:"is_verified"`
}

type UpdateMeCandidateRequest struct {
	updateMe
	University   string   `json:"university" db:"university"`
	Course       int      `json:"course" db:"course"`
	Skills       []string `json:"skills" db:"skills"`
	PortfolioURL string   `json:"portfolio_url" db:"portfolio_url"`
	GithubURL    string   `json:"github_url" db:"github_url"`
}

type VerifyEmployerRequest struct {
	INN string `json:"inn" db:"inn"`
}

type UpdatePrivacyRequest struct {
	IsPrivate bool `json:"is_private"`
}

type UpdateAvatarRequest struct {
	URL string `json:"url"`
}
