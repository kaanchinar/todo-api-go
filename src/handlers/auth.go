package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kaanchinar/todo-app/config"
	"github.com/kaanchinar/todo-app/middleware"
	"github.com/kaanchinar/todo-app/models"
	"github.com/kaanchinar/todo-app/store"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store *store.CompositeStore
	cfg   *config.Config
}

func NewAuthHandler(s *store.CompositeStore, cfg *config.Config) *AuthHandler {
	return &AuthHandler{store: s, cfg: cfg}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration details"
// @Success 201 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "username and password required"})
		return
	}

	if _, exists := h.store.GetUserByUsername(r.Context(), req.Username); exists {
		writeJSON(w, http.StatusConflict, models.ErrorResponse{Error: "username already exists"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to hash password"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashed),
	}
	if err := h.store.CreateUser(r.Context(), user); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create user"})
		return
	}

	writeJSON(w, http.StatusCreated, models.MessageResponse{Message: "user registered successfully"})
}

// Login godoc
// @Summary Login with credentials
// @Description Authenticate with username and password, returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
		return
	}

	user, exists := h.store.GetUserByUsername(r.Context(), req.Username)
	if !exists {
		writeJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		writeJSON(w, http.StatusUnauthorized, models.ErrorResponse{Error: "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(user.ID),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to generate token"})
		return
	}

	writeJSON(w, http.StatusOK, models.LoginResponse{Token: tokenString})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Ensure interface is used
var _ = middleware.UserIDFromContext
