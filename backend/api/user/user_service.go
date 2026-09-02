package user

import (
	"context"
	"errors"
	"server/helper"
	"time"
)

type UserService interface {
	CreateAccount(ctx context.Context, email, password string) error
	VerifyAccount(ctx context.Context, token string) error
	Login(ctx context.Context, email, password string) (string, error)
	ForgetPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, password string) error
}

type UserServiceImpl struct {
	UserRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return &UserServiceImpl{UserRepository: userRepository}
}

func (s *UserServiceImpl) CreateAccount(ctx context.Context, email, password string) error {
	hashPassword, err := helper.HashPassword(password)
	if err != nil {
		return err
	}

	verificationToken := helper.CreateRandomToken()

	userId, err := s.UserRepository.Create(ctx, email, hashPassword, verificationToken)
	if err != nil {
		return err
	}

	if err := helper.SendVerificationEmail(email, verificationToken); err != nil {
		// Roll back user creation if email fails
		s.UserRepository.Delete(ctx, userId)
		return helper.ErrFailedToSendVerificationEmail
	}

	return nil
}

func (s *UserServiceImpl) VerifyAccount(ctx context.Context, token string) error {
	if err := s.UserRepository.SetVerified(ctx, token); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			return "", helper.ErrInvalidEmailOrPassword
		}
		return "", err
	}

	if !helper.CheckPasswordHash(password, user.Password) {
		return "", helper.ErrInvalidEmailOrPassword
	}

	if !user.IsVerified {
		if err := helper.SendVerificationEmail(email, user.VerificationToken); err != nil {
			return "", helper.ErrFailedToSendVerificationEmail
		}
		return "", helper.ErrUserNotVerified
	}

	accessToken, err := helper.CreateAccessToken(user.Id)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *UserServiceImpl) ForgetPassword(ctx context.Context, email string) error {
	resetToken := helper.CreateRandomToken()
	expiry := time.Now().Add(15 * time.Minute)

	if err := s.UserRepository.SetResetToken(ctx, email, resetToken, expiry); err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			return nil
		}
		return err
	}

	if err := helper.SendPasswordResetEmail(email, resetToken); err != nil {
		return helper.ErrFailedToSendVerificationEmail
	}

	return nil
}

func (s *UserServiceImpl) ResetPassword(ctx context.Context, token, password string) error {
	hashPassword, err := helper.HashPassword(password)
	if err != nil {
		return err
	}

	if err := s.UserRepository.ResetPassword(ctx, token, hashPassword); err != nil {
		return err
	}

	return nil
}
