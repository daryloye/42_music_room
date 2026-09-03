package main

import (
	"log"
	"net/http"
	"os"
	"server/api/user"
	"server/prisma"
	"server/router"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	log.Println("Server starting on port", os.Getenv("BACKEND_PORT"))

	db, err := prisma.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to DB:", err)
	}
	defer db.Prisma.Disconnect()

	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userController := user.NewUserController(userService)

	routes := router.NewRouter(userController)

	server := &http.Server{
		Addr:           ":" + os.Getenv("BACKEND_PORT"),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        routes,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Could not start server:", err)
	}
}
