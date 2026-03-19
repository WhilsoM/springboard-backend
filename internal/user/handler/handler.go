package handler

import (
	"log"
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
	val := r.Context().Value(middleware.UserIDKey)
	userID, ok := val.(string)
	if !ok {
		lib.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: user_id not found in context")
		return
	}

	log.Printf("GetMe user id: %s", userID)
	user, err := h.userService.GetMe(r.Context(), userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			lib.WriteErrorJSON(w, http.StatusNotFound, "user not found")
			return
		}
		lib.WriteErrorJSON(w, http.StatusInternalServerError, "failed to get user: "+err.Error())
		return
	}

	lib.WriteJSON(w, http.StatusOK, user)
}
