package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/api/user"
	"server/config"
	"server/controller"
	"server/repository"
	"server/router"
	"server/service"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment variables: ", err)
	}

	fmt.Println("Server starting on port " + os.Getenv("PORT"))

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to DB: ", err)
	}
	defer db.Prisma.Disconnect()

	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository)
	postController := controller.NewPostController(postService)

	userRepository := user.NewUserRepository(db)
	userService := user.NewUserService(userRepository)
	userController := user.NewUserController(userService)

	routes := router.NewRouter(postController, userController)

	server := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        routes,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Could not start server: ", err)
	}
}
