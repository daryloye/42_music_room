package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserController struct {
	UserService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{UserService: userService}
}

type CreateAccountRequest struct {
	Email    string `validate:"required min=1,max=100" json:"email"`
	Password string `validate:"required" json:"password"`
}

type CreateAccountResponse struct {
	Id string `json:"id"`
}

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

	userId, err := c.UserService.CreateAccount(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := CreateAccountResponse{
		Id: userId,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to write response: ", err)
	}
}
