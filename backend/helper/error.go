package helper

import "errors"

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	ErrUserNotFound                  = errors.New("User not found")
	ErrEmailAlreadyExists            = errors.New("Email already exists")
	ErrFailedToSendVerificationEmail = errors.New("Failed to send verification email")
)
