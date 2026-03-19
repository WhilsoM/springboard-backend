package handler

import (
	"net/http"
	"springboard/internal/lib"
	"springboard/internal/middleware"
	"springboard/internal/user/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux, authMW func(http.Handler) http.Handler) {
	mux.Handle("GET /users/me", authMW(http.HandlerFunc(h.GetMe)))
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get data from context
	userID, okID := ctx.Value(middleware.UserIDKey).(string)
	userRole, okRole := ctx.Value(middleware.UserRoleKey).(string)

	if !okID || !okRole {
		lib.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: identity missing in context")
		return
	}

	user, err := h.userService.GetMe(ctx, userID, lib.UserRole(userRole))
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusNotFound, "User profile not found")
		return
	}

	lib.WriteJSON(w, http.StatusOK, user)
}
