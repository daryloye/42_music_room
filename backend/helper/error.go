package helper

import "errors"

var (
	ErrUserNotFound               = errors.New("User not found")
	ErrUserEmailAlreadyExists     = errors.New("Email already exists")
	ErrUserFailedToSendEmail      = errors.New("Failed to send email")
	ErrUserInvalidEmailOrPassword = errors.New("Invalid email or password")
	ErrUserNotVerified            = errors.New("Account is not verified. A new verification email has been sent")
	ErrUserInvalidOrExpiredToken  = errors.New("Invalid or expired token")
)
