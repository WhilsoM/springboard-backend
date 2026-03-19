package lib

type UserRole string

const (
	RoleStudent  UserRole = "applicant"
	RoleEmployer UserRole = "employer"
	RoleCurator  UserRole = "curator"
)
