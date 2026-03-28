package middleware

import (
	"net/http"
	"springboard/internal/lib"
)

// check user role
func RequireRole(allowedRoles ...lib.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleStr, ok := r.Context().Value(UserRoleKey).(string)
			if !ok {
				lib.WriteErrorJSON(w, http.StatusUnauthorized, "role missing in context")
				return
			}

			userRole := lib.UserRole(roleStr)
			hasAccess := false
			for _, allowed := range allowedRoles {
				if userRole == allowed {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				lib.WriteErrorJSON(w, http.StatusForbidden, "access denied: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
