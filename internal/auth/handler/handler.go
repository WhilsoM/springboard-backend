package handler

import (
	"encoding/json"
	"net/http"
	"springboard/internal/auth/dto"
	"springboard/internal/auth/service"
	"springboard/internal/lib"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: s,
	}
}

// register all routes
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/refresh", h.RefreshToken)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.authService.Register(r.Context(), req)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	lib.WriteJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.authService.Login(r.Context(), req)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusUnauthorized, "not correct credentials")
		return
	}

	lib.WriteJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.authService.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	lib.WriteJSON(w, http.StatusOK, response)
}
