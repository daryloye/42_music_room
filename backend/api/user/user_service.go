package user

import (
	"context"
	"server/helper"
)

type UserService interface {
	CreateAccount(ctx context.Context, email, password string) (string, error)
	VerifyAccount(ctx context.Context, token string) error
}

type UserServiceImpl struct {
	UserRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return &UserServiceImpl{UserRepository: userRepository}
}

func (s *UserServiceImpl) CreateAccount(ctx context.Context, email, password string) (string, error) {
	hashPassword, err := helper.HashPassword(password)
	if err != nil {
		return "", err
	}

	verificationToken := helper.CreateRandomToken()

	userId, err := s.UserRepository.Create(ctx, email, hashPassword, verificationToken)
	if err != nil {
		return "", err
	}

	if err := helper.SendVerificationEmail(email, verificationToken); err != nil {
		// Roll back user creation if email fails
		s.UserRepository.Delete(ctx, userId)
		return "", helper.ErrFailedToSendVerificationEmail
	}

	return userId, nil
}

func (s *UserServiceImpl) VerifyAccount(ctx context.Context, token string) error {
	if err := s.UserRepository.SetVerified(ctx, token); err != nil {
		return err
	}

	return nil
}
