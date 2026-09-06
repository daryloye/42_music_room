package user

type User struct {
	Id                string
	Email             string
	Password          string
	DisplayName       string
	VerificationToken string
	IsVerified        bool
}
