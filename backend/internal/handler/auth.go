package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"goodcity/backend/internal/domain"
	"goodcity/backend/internal/service"
)

type AuthHandler struct {
	auth      *service.AuthService
	uploadDir string
}

func NewAuthHandler(auth *service.AuthService, uploadDir string) *AuthHandler {
	return &AuthHandler{auth: auth, uploadDir: uploadDir}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		SocialRole string `json:"social_role"`
		Name       string `json:"name"`
		City       string `json:"city"`
		District   string `json:"district"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	nullableStr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	user, token, err := h.auth.Register(r.Context(), req.Email, req.Password, domain.SocialRole(req.SocialRole),
		nullableStr(req.Name), nullableStr(req.City), nullableStr(req.District))
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			writeError(w, http.StatusConflict, "email already taken")
			return
		}
		writeError(w, http.StatusInternalServerError, "registration failed")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		writeError(w, http.StatusInternalServerError, "login failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	currentUser := UserFromContext(r.Context())
	user, err := h.auth.GetByID(r.Context(), currentUser.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get user")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	currentUser := UserFromContext(r.Context())

	var req struct {
		Name          string `json:"name"`
		Phone         string `json:"phone"`
		City          string `json:"city"`
		District      string `json:"district"`
		SocialRole    string `json:"social_role"`
		IncomeBracket string `json:"income_bracket"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	nullable := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	socialRole := domain.SocialRole(req.SocialRole)
	if socialRole == "" {
		socialRole = domain.SocialResident
	}
	incomeBracket := domain.IncomeBracket(req.IncomeBracket)
	if incomeBracket == "" {
		incomeBracket = domain.IncomeMiddle
	}

	user, err := h.auth.UpdateMe(r.Context(), currentUser.ID,
		nullable(req.Name), nullable(req.Phone), nullable(req.City), nullable(req.District),
		socialRole, incomeBracket,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) SubmitVerification(w http.ResponseWriter, r *http.Request) {
	currentUser := UserFromContext(r.Context())

	if currentUser.Role == domain.RoleModerator {
		writeError(w, http.StatusBadRequest, "moderators do not need role verification")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		writeError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		writeError(w, http.StatusBadRequest, "document file is required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dir := filepath.Join(h.uploadDir, "verifications")
	if err := os.MkdirAll(dir, 0755); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create upload directory")
		return
	}

	dst, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to write file")
		return
	}

	docURL := "/uploads/verifications/" + filename
	if err := h.auth.SubmitVerification(r.Context(), currentUser.ID, docURL); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to submit verification")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"doc_url": docURL})
}
