package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/smetanamolokovich/veylo/internal/application/auth"
	"github.com/smetanamolokovich/veylo/internal/domain/refreshtoken"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

type AuthHandler struct {
	registerUseCase     *auth.RegisterUseCase
	loginUseCase        *auth.LoginUseCase
	refreshTokenUseCase *auth.RefreshTokenUseCase
	signupUseCase       *auth.SignupUseCase
}

func NewAuthHandler(registerUC *auth.RegisterUseCase, loginUC *auth.LoginUseCase, refreshUC *auth.RefreshTokenUseCase, signupUC *auth.SignupUseCase) *AuthHandler {
	return &AuthHandler{
		registerUseCase:     registerUC,
		loginUseCase:        loginUC,
		refreshTokenUseCase: refreshUC,
		signupUseCase:       signupUC,
	}
}

type registerRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	OrganizationID string `json:"organization_id"`
	FullName       string `json:"full_name"`
	Role           string `json:"role"`
}

type loginRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	OrganizationID string `json:"organization_id"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.registerUseCase.Execute(r.Context(), auth.RegisterRequest{
		Email:          req.Email,
		Password:       req.Password,
		OrganizationID: req.OrganizationID,
		FullName:       req.FullName,
		Role:           user.Role(req.Role),
	})
	if err != nil {
		if errors.Is(err, user.ErrAlreadyExists) {
			writeError(w, http.StatusConflict, "email already in use")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.loginUseCase.Execute(r.Context(), auth.LoginRequest{
		Email:          req.Email,
		Password:       req.Password,
		OrganizationID: req.OrganizationID,
	})
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound) || errors.Is(err, user.ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid email or password")
			return
		case errors.Is(err, user.ErrBlocked):
			writeError(w, http.StatusForbidden, "account is blocked")
			return
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, resp)
}

type refreshRequest struct {
	RefreshToken   string `json:"refresh_token"`
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrgName  string `json:"org_name"`
		Vertical string `json:"vertical"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.signupUseCase.Execute(r.Context(), auth.SignupRequest{
		OrgName:  req.OrgName,
		Vertical: req.Vertical,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.refreshTokenUseCase.Execute(r.Context(), auth.RefreshRequest{
		RefreshToken:   req.RefreshToken,
		UserID:         req.UserID,
		OrganizationID: req.OrganizationID,
	})
	if err != nil {
		switch {
		case errors.Is(err, refreshtoken.ErrNotFound) || errors.Is(err, refreshtoken.ErrInvalidRefreshToken):
			writeError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		case errors.Is(err, refreshtoken.ErrExpiredRefreshToken):
			writeError(w, http.StatusUnauthorized, "refresh token expired")
			return
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, resp)
}
