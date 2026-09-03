package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"server/helper"
	"time"

	"github.com/julienschmidt/httprouter"
)

type UserController struct {
	UserService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{UserService: userService}
}

// @Summary	Create account
// @Param request body CreateAccountRequest true "Email and password"
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

	if request.Email == "" || request.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	err := c.UserService.CreateAccount(r.Context(), request.Email, request.Password)
	if err != nil {
		log.Println("Failed to create account:", err)
		switch {
		case errors.Is(err, helper.ErrEmailAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, helper.ErrFailedToSendVerificationEmail):
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Registration successful! Please check your email to verify your account",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

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
		case errors.Is(err, helper.ErrInvalidOrExpiredToken):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Email verified successfully! You can now log in",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

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

	accessToken, err := c.UserService.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		log.Println("Failed to login:", err)
		switch {
		case errors.Is(err, helper.ErrInvalidEmailOrPassword):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, helper.ErrUserNotVerified):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Login successful",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		// Secure: true,						// TODO for HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write response:", err)
	}
}

// @Summary Logout
// @Success 200
// @Failure 500
// @Router /api/auth/logout [post]
func (c *UserController) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	response := map[string]string{
		"message": "Logged out",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	cookie := &http.Cookie{
		Name:     "access_token",
		MaxAge:   -1,
		HttpOnly: true,
		// Secure: true, 						// TODO for HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

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

	response := map[string]string{
		"message": "A password reset link has been sent",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

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
		case errors.Is(err, helper.ErrInvalidOrExpiredToken):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Password updated",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write response:", err)
	}
}
