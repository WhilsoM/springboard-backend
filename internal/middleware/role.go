package middleware

import (
	"net/http"
	"springboard/internal/lib"
)

// check user role
func RoleMiddleware(allowedRoles ...lib.UserRole) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("role").(lib.UserRole)
			if !ok {
				lib.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// check if user role is allowed
			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			lib.WriteErrorJSON(w, http.StatusForbidden, "You don't have enough permissions")
		})
	}
}
