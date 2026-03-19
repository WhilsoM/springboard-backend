package lib

import "golang.org/x/crypto/bcrypt"

type User struct {
	PasswordHash string `json:"password_hash" db:"password_hash"`
	ID           string `json:"id" db:"id"`
	Role         string `json:"role" db:"role"`
	Email        string `json:"email" db:"email"`
	FullName     string `json:"full_name" db:"full_name"`
}

func GenerateHashByPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
