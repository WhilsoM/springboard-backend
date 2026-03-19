package dto

import "springboard/internal/lib"

type UserResponse struct {
	ID       string       `json:"id"`
	Role     lib.UserRole `json:"role"`
	Email    string       `json:"email"`
	FullName string       `json:"full_name"`
}
