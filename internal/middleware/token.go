package middleware

import (
	"context"
	"net/http"
	"springboard/internal/lib"
	"strings"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UserRoleKey  contextKey = "user_role"
)

func CheckTokenMiddleware(jwtManager *lib.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				lib.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: invalid header format")
				return
			}

			claims, err := jwtManager.ValidateToken(parts[1])
			if err != nil {
				lib.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: invalid token")
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
