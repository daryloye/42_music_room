package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	EnvAppHostname      string
	EnvBackendPort      string
	EnvFrontendPort     string
	EnvEmailUser        string
	EnvEmailPassword    string
	EnvJwtAccessSecret  string
	EnvJwtAccessExpiry  time.Duration
	EnvJwtRefreshExpiry time.Duration
}

func NewConfig() (*Config, error) {
	godotenv.Load()

	accessExpiry, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXPIRY"))
	if err != nil {
		return &Config{}, err
	}

	refreshExpiry, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRY"))
	if err != nil {
		return &Config{}, err
	}

	return &Config{
		EnvAppHostname:  os.Getenv("APP_HOSTNAME"),
		EnvBackendPort:  os.Getenv("BACKEND_PORT"),
		EnvFrontendPort: os.Getenv("FRONTEND_PORT"),

		EnvEmailUser:     os.Getenv("EMAIL_USER"),
		EnvEmailPassword: os.Getenv("EMAIL_PASSWORD"),

		EnvJwtAccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		EnvJwtAccessExpiry:  accessExpiry,
		EnvJwtRefreshExpiry: refreshExpiry,
	}, nil
}
