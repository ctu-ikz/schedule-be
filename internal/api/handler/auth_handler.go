package handler

import (
	"encoding/json"
	"errors"
	"github.com/ctu-ikz/schedule-be/internal/api/dto"
	"github.com/ctu-ikz/schedule-be/internal/service"
	"github.com/ctu-ikz/schedule-be/internal/util"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
	"time"
)

var validate = validator.New()

type AuthHandler struct {
	service *service.AuthService
}

func NewUserHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	res := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	refreshToken, accessToken, err := h.service.Login(r.Context(), req.Username, req.Password, r.UserAgent(), util.GetIP(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := map[string]string{
		"token_access": accessToken,
	}

	secure := os.Getenv("APP_ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth/refresh",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		}
		http.Error(w, "Error reading refresh token", http.StatusInternalServerError)
		return
	}

	refreshTokenValue := cookie.Value

	refreshToken, accessToken, err := h.service.Refresh(r.Context(), refreshTokenValue, r.UserAgent(), util.GetIP(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	secure := os.Getenv("APP_ENV") == "production"

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Path:     "/auth/refresh",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
	})

	resp := map[string]string{
		"token_access": accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
