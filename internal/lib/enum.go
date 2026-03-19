package lib

type UserRole string

const (
	RoleStudent  UserRole = "applicant"
	RoleEmployer UserRole = "employer"
	RoleCurator  UserRole = "curator"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleStudent, RoleEmployer, RoleCurator:
		return true
	}
	return false
}
