package user

import (
	"context"
	"errors"
	"server/api/refreshtoken"
	"server/config"
	"server/helper"
	"time"
)

type UserService interface {
	CreateAccount(ctx context.Context, email, password, displayName string) error
	VerifyAccount(ctx context.Context, token string) error
	Login(ctx context.Context, email, password string) (string, string, error)
	ForgetPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, password string) error
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	GetProfile(ctx context.Context, id string) (User, error)
	UpdateProfile(ctx context.Context, id string, request UpdateProfileRequest) (User, error)
	UpdatePassword(ctx context.Context, id, password string) error
}

type UserServiceImpl struct {
	Cfg                 *config.Config
	UserRepository      UserRepository
	RefreshTokenService refreshtoken.RefreshTokenService
}

func NewUserService(cfg *config.Config, userRepository UserRepository, refreshTokenService refreshtoken.RefreshTokenService) UserService {
	return &UserServiceImpl{
		Cfg:                 cfg,
		UserRepository:      userRepository,
		RefreshTokenService: refreshTokenService,
	}
}

func (s *UserServiceImpl) CreateAccount(ctx context.Context, email, password, displayName string) error {
	hashPassword, err := helper.HashPassword(password)
	if err != nil {
		return err
	}

	verificationToken := helper.CreateRandomToken()

	userId, err := s.UserRepository.Create(ctx, email, hashPassword, displayName, verificationToken)
	if err != nil {
		return err
	}

	if err := helper.SendVerificationEmail(s.Cfg, email, verificationToken); err != nil {
		// Roll back user creation if email fails
		s.UserRepository.Delete(ctx, userId)
		return helper.ErrUserFailedToSendEmail
	}

	return nil
}

func (s *UserServiceImpl) VerifyAccount(ctx context.Context, token string) error {
	if err := s.UserRepository.SetVerified(ctx, token); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, helper.ErrUserNotFound) {
			return "", "", helper.ErrUserInvalidEmailOrPassword
		}
		return "", "", err
	}

	if !helper.CheckPasswordHash(password, user.Password) {
		return "", "", helper.ErrUserInvalidEmailOrPassword
	}

	if !user.IsVerified {
		if err := helper.SendVerificationEmail(s.Cfg, email, user.VerificationToken); err != nil {
			return "", "", helper.ErrUserFailedToSendEmail
		}
		return "", "", helper.ErrUserNotVerified
	}

	accessToken, err := helper.CreateJWT(user.Id, s.Cfg.EnvJwtAccessSecret, s.Cfg.EnvJwtAccessExpiry)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.RefreshTokenService.Create(ctx, user.Id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
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

	if err := helper.SendPasswordResetEmail(s.Cfg, email, resetToken); err != nil {
		return helper.ErrUserFailedToSendEmail
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

func (s *UserServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	if err := s.RefreshTokenService.Delete(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) RefreshToken(ctx context.Context, oldToken string) (string, string, error) {
	userId, newRefreshToken, err := s.RefreshTokenService.Update(ctx, oldToken)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := helper.CreateJWT(userId, s.Cfg.EnvJwtAccessSecret, s.Cfg.EnvJwtAccessExpiry)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, id string) (User, error) {
	user, err := s.UserRepository.FindById(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, id string, request UpdateProfileRequest) (User, error) {
	user, err := s.UserRepository.Update(ctx, id, request)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *UserServiceImpl) UpdatePassword(ctx context.Context, id, password string) error {
	hashPassword, err := helper.HashPassword(password)
	if err != nil {
		return err
	}

	if err := s.UserRepository.UpdatePassword(ctx, id, hashPassword); err != nil {
		return err
	}

	return nil
}
