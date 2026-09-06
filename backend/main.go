package main

import (
	"log"
	"net/http"
	"server/api/refreshtoken"
	"server/api/user"
	"server/config"
	"server/prisma"
	"server/router"
	"time"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Could not load env:", err)
	}

	db, err := prisma.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to DB:", err)
	}
	defer db.Prisma.Disconnect()

	refreshTokenRepository := refreshtoken.NewRefreshTokenRepository(db)
	userRepository := user.NewUserRepository(db)

	refreshTokenService := refreshtoken.NewRefreshTokenService(cfg, refreshTokenRepository)
	userService := user.NewUserService(cfg, userRepository, refreshTokenService)

	userController := user.NewUserController(cfg, userService)

	routes := router.NewRouter(userController)

	server := &http.Server{
		Addr:           ":" + cfg.EnvBackendPort,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        routes,
	}

	log.Println("Server starting on port", cfg.EnvBackendPort)

	if err = server.ListenAndServe(); err != nil {
		log.Fatal("Could not start server:", err)
	}
}
