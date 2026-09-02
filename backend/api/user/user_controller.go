package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"server/helper"

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

	userId, err := c.UserService.CreateAccount(r.Context(), request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, helper.ErrEmailAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, helper.ErrFailedToSendVerificationEmail):
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			log.Println("Failed to create account:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{
		"message": "Registration successful! Please check your email to verify your account",
		"id":      userId,
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (c *UserController) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	// TODO
}
