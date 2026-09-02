package router

import (
	"encoding/json"
	"log"
	"net/http"
	"server/api/user"
	"server/controller"
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

func NewRouter(postController *controller.PostController, userController *user.UserController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/health", healthCheckHandler)

	router.GET("/api/post", postController.FindAll)
	router.GET("/api/post/:postId", postController.FindById)
	router.POST("/api/post", postController.Create)
	router.PATCH("/api/post/:postId", postController.Update)
	router.DELETE("/api/post/:postId", postController.Delete)

	router.POST("/api/auth/signup", userController.CreateAccount)
	router.GET("/api/auth/verify", userController.VerifyAccount)
	router.POST("/api/auth/login", userController.Login)
	router.POST("/api/auth/logout", userController.Logout)
	router.POST("/api/auth/forget-password", userController.ForgetPassword)
	router.POST("/api/auth/reset-password", userController.ResetPassword)

	return router
}
