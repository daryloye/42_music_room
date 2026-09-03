package helper

import "errors"

var (
	ErrUserNotFound                  = errors.New("User not found")
	ErrEmailAlreadyExists            = errors.New("Email already exists")
	ErrFailedToSendVerificationEmail = errors.New("Failed to send verification email")
	ErrInvalidEmailOrPassword        = errors.New("Invalid email or password")
	ErrUserNotVerified               = errors.New("Account is not verified. A new verification email has been sent")
	ErrInvalidOrExpiredToken         = errors.New("Invalid or expired token")
)
