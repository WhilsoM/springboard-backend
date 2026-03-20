package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"springboard/internal/lib"
	"springboard/internal/middleware"
	"springboard/internal/opportunity/dto"
	"springboard/internal/opportunity/service"
)

type OpportunityHandler struct {
	service service.OpportunityService
}

func NewOpportunityHandler(s service.OpportunityService) *OpportunityHandler {
	return &OpportunityHandler{service: s}
}

func (h *OpportunityHandler) RegisterRoutes(mux *http.ServeMux, authMW func(http.Handler) http.Handler) {
	mux.Handle("GET /opportunities", http.HandlerFunc(h.List))
	mux.Handle("GET /opportunities/{id}", http.HandlerFunc(h.GetByID))

	mux.Handle("POST /opportunities", authMW(http.HandlerFunc(h.Create)))
	mux.Handle("PUT /opportunities/{id}", authMW(http.HandlerFunc(h.Update)))
	mux.Handle("DELETE /opportunities/{id}", authMW(http.HandlerFunc(h.Delete)))
	mux.Handle("GET /employer/opportunities", authMW(http.HandlerFunc(h.GetMy)))
}

func (h *OpportunityHandler) List(w http.ResponseWriter, r *http.Request) {
	filters := dto.SearchFilters{
		Type:   r.URL.Query().Get("type"),
		Format: r.URL.Query().Get("format"),
		Search: r.URL.Query().Get("search"),
	}
	log.Println("List get filters", filters)
	res, err := h.service.GetAll(r.Context(), filters)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("result in the handler", res)
	lib.WriteJSON(w, http.StatusOK, map[string][]dto.OpportunityResponse{"opportunities": res})
}

func (h *OpportunityHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, isOk := r.Context().Value(middleware.UserIDKey).(string)
	if !isOk {
		lib.WriteErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateOpportunityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	res, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusForbidden, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusCreated, res)
}

func (h *OpportunityHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, err := h.service.GetOne(r.Context(), id)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusNotFound, "not found")
		return
	}
	lib.WriteJSON(w, http.StatusOK, res)
}

func (h *OpportunityHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req dto.CreateOpportunityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.service.Update(r.Context(), id, userID, req); err != nil {
		lib.WriteErrorJSON(w, http.StatusForbidden, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *OpportunityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.Context().Value(middleware.UserIDKey).(string)

	if err := h.service.Delete(r.Context(), id, userID); err != nil {
		lib.WriteErrorJSON(w, http.StatusForbidden, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *OpportunityHandler) GetMy(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	res, err := h.service.GetEmployerOwn(r.Context(), userID)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, res)
}
