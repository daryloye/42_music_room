package user

type User struct {
	Id                string `json:"id,omitempty"`
	Email             string `json:"email,omitempty"`
	Password          string `json:"password,omitempty"`
	DisplayName       string `json:"display_name,omitempty"`
	VerificationToken string `json:"verification_token,omitempty"`
	IsVerified        bool   `json:"is_verified,omitempty"`
}
