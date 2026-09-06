package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"server/config"
	"server/helper"
	"time"

	"github.com/julienschmidt/httprouter"
)

type UserController struct {
	Cfg         *config.Config
	UserService UserService
}

func NewUserController(cfg *config.Config, userService UserService) *UserController {
	return &UserController{Cfg: cfg, UserService: userService}
}

// @Summary	Create account
// @Param request body CreateAccountRequest true "Email, password, display name"
// @Success 201
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /api/auth/signup [post]
func (c *UserController) CreateAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var request CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if request.Email == "" || request.Password == "" || request.DisplayName == "" {
		http.Error(w, "Email, password and display name are required", http.StatusBadRequest)
		return
	}

	err := c.UserService.CreateAccount(r.Context(), request.Email, request.Password, request.DisplayName)
	if err != nil {
		log.Println("Failed to create account:", err)
		switch {
		case errors.Is(err, helper.ErrUserEmailAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, helper.ErrUserFailedToSendEmail):
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := map[string]string{
		"message": "Registration successful! Please check your email to verify your account",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write response:", err)
	}
}

// @Summary	Verify account
// @Param token query string true "Verification token"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/auth/verify [get]
func (c *UserController) VerifyAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var request VerifyAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if request.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	if err := c.UserService.VerifyAccount(r.Context(), request.Token); err != nil {
		log.Println("Failed to verify account:", err)
		switch {
		case errors.Is(err, helper.ErrUserInvalidOrExpiredToken):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Email verified successfully! You can now log in",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write response:", err)
	}
}

// @Summary Login
// @Params request body LoginRequest true "Email and password"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 403
// @Failure 500
// @Router /api/auth/login [post]
func (c *UserController) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if request.Email == "" || request.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := c.UserService.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		log.Println("Failed to login:", err)
		switch {
		case errors.Is(err, helper.ErrUserInvalidEmailOrPassword):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, helper.ErrUserNotVerified):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(c.Cfg.EnvJwtAccessExpiry),
		MaxAge:   int(c.Cfg.EnvJwtAccessExpiry.Seconds()),
		HttpOnly: true,
		// Secure: true,						// TODO for HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(c.Cfg.EnvJwtRefreshExpiry),
		MaxAge:   int(c.Cfg.EnvJwtRefreshExpiry.Seconds()),
		HttpOnly: true,
		// Secure: true,						// TODO for HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Login successful",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write response:", err)
	}
}

// @Summary Logout
// @Success 200
// @Failure 500
// @Router /api/auth/logout [post]
func (c *UserController) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	refreshToken, _ := r.Cookie("refresh_token")
	if refreshToken != nil && refreshToken.Value != "" {
		c.UserService.Logout(r.Context(), refreshToken.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		// Secure: true, 						// TODO for HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		// Secure: true, 						// TODO for HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Logged out",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write response:", err)
	}
}

// @Summary	Forget password
// @Param request body ForgetPasswordRequest true "Email"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/auth/forget-password [post]
func (c *UserController) ForgetPassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var request ForgetPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if request.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if err := c.UserService.ForgetPassword(r.Context(), request.Email); err != nil {
		log.Println("Failed to reset password:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "A password reset link has been sent",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write response:", err)
	}
}

// @Summary	Reset password
// @Param request body ResetPasswordRequest true "Token and password"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/auth/reset-password [post]
func (c *UserController) ResetPassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var request ResetPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if request.Token == "" || request.Password == "" {
		http.Error(w, "Token and password are required", http.StatusBadRequest)
		return
	}

	if err := c.UserService.ResetPassword(r.Context(), request.Token, request.Password); err != nil {
		log.Println("Failed to reset password:", err)
		switch {
		case errors.Is(err, helper.ErrUserInvalidOrExpiredToken):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Password updated",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write response:", err)
	}
}

// @Summary Refresh Token
// @Success 200
// @Failure 401
// @Failure 500
// @Router /api/auth/refresh [post]
func (c *UserController) RefreshToken(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil || refreshToken.Value == "" {
		http.Error(w, helper.ErrUserInvalidOrExpiredToken.Error(), http.StatusUnauthorized)
		return
	}

	newAccessToken, newRefreshToken, err := c.UserService.RefreshToken(r.Context(), refreshToken.Value)
	if err != nil {
		log.Println("Failed to refresh token:", err)
		switch {
		case errors.Is(err, helper.ErrUserInvalidOrExpiredToken):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Path:     "/",
		Expires:  time.Now().Add(c.Cfg.EnvJwtAccessExpiry),
		MaxAge:   int(c.Cfg.EnvJwtAccessExpiry.Seconds()),
		HttpOnly: true,
		// Secure: true,						// TODO for HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(c.Cfg.EnvJwtRefreshExpiry),
		MaxAge:   int(c.Cfg.EnvJwtRefreshExpiry.Seconds()),
		HttpOnly: true,
		// Secure: true,						// TODO for HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Token refreshed",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write response:", err)
	}
}
