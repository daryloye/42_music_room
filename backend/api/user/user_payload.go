package user

type CreateAccountRequest struct {
	Email    string `validate:"required,min=1,max=100" json:"email"`
	Password string `validate:"required" json:"password"`
}

type VerifyAccountRequest struct {
	Token string `validate:"required" json:"token"`
}

type LoginRequest struct {
	Email    string `validate:"required,min=1,max=100" json:"email"`
	Password string `validate:"required" json:"password"`
}
