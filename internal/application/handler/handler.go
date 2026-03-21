package handler

import (
	"encoding/json"
	"net/http"
	"springboard/internal/application/dto"
	"springboard/internal/application/service"
	"springboard/internal/lib"
	"springboard/internal/middleware"
	"strconv"
)

type ApplicationHandler struct {
	service service.ApplicationService
}

func NewApplicationHandler(s service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: s}
}

func (h *ApplicationHandler) RegisterRoutes(mux *http.ServeMux, authMW func(http.Handler) http.Handler) {
	mux.Handle("POST /opportunities/{id}/apply", authMW(http.HandlerFunc(h.Apply)))
	mux.Handle("GET /employer/opportunities/{id}/applications", authMW(http.HandlerFunc(h.GetForEmployer)))
	mux.Handle("PATCH /employer/applications/{app_id}/status", authMW(http.HandlerFunc(h.UpdateStatus)))
}

func (h *ApplicationHandler) Apply(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid := r.Context().Value(middleware.UserIDKey).(string)
	role := r.Context().Value(middleware.UserRoleKey).(string)

	var req dto.ApplyRequest
	json.NewDecoder(r.Body).Decode(&req)

	res, err := h.service.Apply(r.Context(), uid, role, id, req)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusForbidden, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusCreated, res)
}

func (h *ApplicationHandler) GetForEmployer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	uid := r.Context().Value(middleware.UserIDKey).(string)
	role := r.Context().Value(middleware.UserRoleKey).(string)

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 10
	}

	res, err := h.service.GetForEmployer(r.Context(), uid, role, id, dto.PaginationFilters{Limit: limit, Offset: offset})
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusForbidden, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string]any{"applications": res, "limit": limit, "offset": offset})
}

func (h *ApplicationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("app_id")
	uid := r.Context().Value(middleware.UserIDKey).(string)
	role := r.Context().Value(middleware.UserRoleKey).(string)

	var req dto.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.service.UpdateStatus(r.Context(), uid, role, appID, req); err != nil {
		lib.WriteErrorJSON(w, http.StatusForbidden, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string]string{"message": "status updated"})
}
