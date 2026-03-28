package dto

type CreateCuratorRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateTagRequest struct {
	Name string `json:"name"`
}

type Verification struct {
	ID          string `json:"id"`
	EmployerID  string `json:"employer_id"`
	INN         string `json:"inn"`
	CompanyName string `json:"company_name"`
	Status      string `json:"status"`
}

type UpdateVerificationRequest struct {
	EmployerID string `json:"employer_id"`
	Status     string `json:"status"` // 'approved' or 'rejected'
}
