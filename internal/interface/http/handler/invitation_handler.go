package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	appinvitation "github.com/smetanamolokovich/veylo/internal/application/invitation"
	doaminvitation "github.com/smetanamolokovich/veylo/internal/domain/invitation"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type InvitationHandler struct {
	inviteUseCase  *appinvitation.InviteUserUseCase
	getUseCase     *appinvitation.GetInvitationUseCase
	acceptUseCase  *appinvitation.AcceptInvitationUseCase
}

func NewInvitationHandler(
	inviteUC *appinvitation.InviteUserUseCase,
	getUC *appinvitation.GetInvitationUseCase,
	acceptUC *appinvitation.AcceptInvitationUseCase,
) *InvitationHandler {
	return &InvitationHandler{
		inviteUseCase: inviteUC,
		getUseCase:    getUC,
		acceptUseCase: acceptUC,
	}
}

type createInvitationRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *InvitationHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok || orgID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, ok := authmiddleware.UserIDFromCtx(r.Context())
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Role == "" {
		writeError(w, http.StatusBadRequest, "email and role are required")
		return
	}

	resp, err := h.inviteUseCase.Execute(r.Context(), appinvitation.InviteUserRequest{
		OrganizationID: orgID,
		InviterUserID:  userID,
		Email:          req.Email,
		Role:           req.Role,
	})
	if err != nil {
		switch {
		case errors.Is(err, doaminvitation.ErrDuplicate):
			writeError(w, http.StatusConflict, "pending invitation already exists for this email")
		case errors.Is(err, user.ErrAlreadyExists):
			writeError(w, http.StatusConflict, "user already exists in this organization")
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *InvitationHandler) GetByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}

	resp, err := h.getUseCase.Execute(r.Context(), appinvitation.GetInvitationRequest{Token: token})
	if err != nil {
		if errors.Is(err, doaminvitation.ErrNotFound) {
			writeError(w, http.StatusNotFound, "invitation not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

type acceptInvitationRequest struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

func (h *InvitationHandler) Accept(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}

	var req acceptInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FullName == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "full_name and password are required")
		return
	}

	resp, err := h.acceptUseCase.Execute(r.Context(), appinvitation.AcceptInvitationRequest{
		Token:    token,
		FullName: req.FullName,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, doaminvitation.ErrNotFound):
			writeError(w, http.StatusNotFound, "invitation not found")
		case errors.Is(err, doaminvitation.ErrExpired):
			writeError(w, http.StatusGone, "invitation has expired")
		case errors.Is(err, doaminvitation.ErrAlreadyUsed):
			writeError(w, http.StatusConflict, "invitation already used")
		case errors.Is(err, user.ErrAlreadyExists):
			writeError(w, http.StatusConflict, "email already registered")
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}
