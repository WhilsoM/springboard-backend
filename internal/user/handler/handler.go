package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"springboard/internal/lib"
	"springboard/internal/middleware"
	"springboard/internal/user/dto"
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
	mux.Handle("DELETE /users/me", authMW(http.HandlerFunc(h.DeleteMe)))
	mux.Handle("PUT /users/me", authMW(http.HandlerFunc(h.UpdateMe)))
	mux.Handle("POST /users/me/verify", authMW(http.HandlerFunc(h.Verify)))
	mux.Handle("PATCH /users/me/privacy", authMW(http.HandlerFunc(h.UpdatePrivacy)))
	mux.Handle("PATCH /users/me/avatar", authMW(http.HandlerFunc(h.UpdateAvatar)))
	mux.Handle("GET /users/{id}", authMW(http.HandlerFunc(h.GetUserProfile)))
	// NETWORK ROUTES
	mux.Handle("GET /applicants", authMW(http.HandlerFunc(h.SearchApplicants)))
	mux.Handle("POST /network/request/{id}", authMW(http.HandlerFunc(h.SendRequest)))
	mux.Handle("PATCH /network/request/{request_id}", authMW(http.HandlerFunc(h.AcceptRejectRequest)))
	mux.Handle("GET /network/contacts", authMW(http.HandlerFunc(h.GetContacts)))
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
	log.Println("user", user)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusNotFound, "User profile not found")
		return
	}

	lib.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, okID := ctx.Value(middleware.UserIDKey).(string)

	log.Println("user id deleteme:", userID)

	if !okID {
		lib.WriteErrorJSON(w, http.StatusUnauthorized, "Unauthorized: identity missing in context")
		return
	}

	err := h.userService.DeleteMe(ctx, userID)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusNotFound, "User profile not found")
		return
	}

	lib.WriteJSON(w, http.StatusOK, map[string]string{"message": "successfull"})
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(middleware.UserIDKey).(string)
	roleStr, _ := ctx.Value(middleware.UserRoleKey).(string)
	userRole := lib.UserRole(roleStr)

	var result any
	var err error

	switch userRole {
	case lib.RoleStudent:
		var req dto.UpdateMeCandidateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid candidate request body")
			return
		}
		result, err = h.userService.UpdateCandidate(ctx, userID, req)

	case lib.RoleEmployer:
		var req dto.UpdateMeEmployerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid employer request body")
			return
		}
		result, err = h.userService.UpdateEmployer(ctx, userID, req)

	default:
		lib.WriteErrorJSON(w, http.StatusForbidden, "this role cannot update profile")
		return
	}

	if err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, "update failed: "+err.Error())
		return
	}

	lib.WriteJSON(w, http.StatusOK, map[string]any{
		"user": result,
	})
}

func (h *UserHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyEmployerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	role := lib.UserRole(r.Context().Value(middleware.UserRoleKey).(string))

	if err := h.userService.Verify(r.Context(), userID, role, req.INN); err != nil {
		lib.WriteErrorJSON(w, http.StatusForbidden, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string]string{"message": "verification request submitted"})
}

func (h *UserHandler) UpdatePrivacy(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePrivacyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if err := h.userService.SetPrivacy(r.Context(), userID, req.IsPrivate); err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string]bool{"is_private": req.IsPrivate})
}

func (h *UserHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	role := lib.UserRole(r.Context().Value(middleware.UserRoleKey).(string))

	if err := h.userService.SetAvatar(r.Context(), userID, role, req.AvatarURL); err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string]string{"avatar_url": req.AvatarURL})
}

func (h *UserHandler) SearchApplicants(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	applicants, err := h.userService.SearchApplicants(r.Context(), query, 20, 0)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string][]lib.ApplicantUser{"applicants": applicants})
}

func (h *UserHandler) SendRequest(w http.ResponseWriter, r *http.Request) {
	senderID := r.Context().Value(middleware.UserIDKey).(string)
	receiverID := r.PathValue("id")

	if err := h.userService.SendRequest(r.Context(), senderID, receiverID); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusCreated, map[string]string{"message": "request sent"})
}

func (h *UserHandler) AcceptRejectRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	requestID := r.PathValue("request_id")

	var req dto.HandleContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lib.WriteErrorJSON(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.userService.HandleContactRequest(r.Context(), userID, requestID, req.Status); err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string]string{"status": req.Status})
}

func (h *UserHandler) GetContacts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	contacts, err := h.userService.GetMyContacts(r.Context(), userID)
	if err != nil {
		lib.WriteErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	lib.WriteJSON(w, http.StatusOK, map[string][]lib.User{"contacts": contacts})
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	targetID := r.PathValue("id")
	profile, err := h.userService.GetUserProfile(r.Context(), targetID)
	if err != nil {
		if err.Error() == "profile is private" {
			lib.WriteErrorJSON(w, http.StatusForbidden, "This profile is private")
			return
		}
		lib.WriteErrorJSON(w, http.StatusNotFound, "User not found")
		return
	}
	lib.WriteJSON(w, http.StatusOK, profile)
}
