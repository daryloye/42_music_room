package user

import "time"

type User struct {
	Id                string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Email             string
	Password          string
	VerificationToken string
	IsVerified        bool
}
