package handler

import (
	"encoding/json"
	"net/http"
	"springboard/internal/admin/dto"
	"springboard/internal/admin/service"
	"springboard/internal/lib"
	"springboard/internal/middleware"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux, authMW func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /tags", h.GetTags)

	// Кураторские ручки
	adminOnly := middleware.RequireRole(lib.RoleCurator)
	employerOnly := middleware.RequireRole(lib.RoleEmployer)

	mux.Handle("POST /admin/curators", authMW(adminOnly(http.HandlerFunc(h.CreateCurator))))
	mux.Handle("POST /tags", authMW(adminOnly(http.HandlerFunc(h.CreateTag))))
	mux.Handle("PATCH /admin/verifications/{id}/status", authMW(adminOnly(http.HandlerFunc(h.ModerateVerification))))
	mux.Handle("DELETE /admin/opportunities/{id}", authMW(adminOnly(http.HandlerFunc(h.DeleteOpportunity))))

	mux.Handle("POST /employer/verify", authMW(employerOnly(http.HandlerFunc(h.SubmitVerification))))
}

func (h *AdminHandler) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.adminService.GetAllTags(r.Context())
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, tags)
}

func (h *AdminHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	tag, err := h.adminService.CreateTag(r.Context(), req.Name)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusCreated, tag)
}

func (h *AdminHandler) CreateCurator(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCuratorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.adminService.CreateCurator(r.Context(), req); err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusCreated, map[string]string{"message": "curator created"})
}

func (h *AdminHandler) SubmitVerification(w http.ResponseWriter, r *http.Request) {
	// todo: to dto
	var req struct {
		Inn         string `json:"inn"`
		CompanyName string `json:"company_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	ctx := r.Context()
	employerID := ctx.Value(middleware.UserIDKey).(string)

	if err := h.adminService.SubmitVerification(ctx, employerID, req.Inn, req.CompanyName); err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	lib.WriteJSON(w, http.StatusOK, map[string]string{"message": "verification processed"})
}

func (h *AdminHandler) ModerateVerification(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("id")
	var req dto.UpdateVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.adminService.ModerateVerification(r.Context(), requestID, req.EmployerID, req.Status); err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string]string{"message": "status updated"})
}

func (h *AdminHandler) DeleteOpportunity(w http.ResponseWriter, r *http.Request) {
	oppID := r.PathValue("id")
	if err := h.adminService.ForceDeleteOpportunity(r.Context(), oppID); err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string]string{"message": "content deleted by admin"})
}
