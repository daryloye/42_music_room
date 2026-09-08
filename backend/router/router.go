package router

import (
	"encoding/json"
	"log"
	"net/http"
	"server/api/user"
	"server/middleware"
	"time"

	"github.com/julienschmidt/httprouter"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	response := struct {
		Status    string    `json:"status"`
		Timestamp time.Time `json:"timestamp"`
	}{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Failed to write health response:", err)
	}
}

func NewRouter(middleware middleware.Middleware, userController *user.UserController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/health", healthCheckHandler)

	router.POST("/api/auth/signup", userController.CreateAccount)
	router.GET("/api/auth/verify", userController.VerifyAccount)
	router.POST("/api/auth/login", userController.Login)
	router.POST("/api/auth/logout", userController.Logout)
	router.POST("/api/auth/forget-password", userController.ForgetPassword)
	router.POST("/api/auth/reset-password", userController.ResetPassword)
	router.POST("/api/auth/refresh", userController.RefreshToken)

	router.GET("/api/profile/me", middleware.RequireAuth(userController.GetProfile))
	router.PUT("/api/profile/me", middleware.RequireAuth(userController.UpdateProfile))
	router.PUT("/api/profile/update-password", middleware.RequireAuth(userController.UpdatePassword))

	return router
}
