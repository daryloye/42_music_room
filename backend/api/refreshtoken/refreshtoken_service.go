package refreshtoken

import (
	"context"
	"server/config"
	"server/helper"
	"time"
)

type RefreshTokenService interface {
	Create(ctx context.Context, userId string) (string, error)
	Update(ctx context.Context, oldToken string) (string, string, error)
	Delete(ctx context.Context, token string) error
}

type RefreshTokenImpl struct {
	Cfg                    *config.Config
	RefreshTokenRepository RefreshTokenRepository
}

func NewRefreshTokenService(cfg *config.Config, refreshTokenRepository RefreshTokenRepository) RefreshTokenService {
	return &RefreshTokenImpl{
		Cfg:                    cfg,
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (s *RefreshTokenImpl) Create(ctx context.Context, userId string) (string, error) {
	token := helper.CreateRandomToken()

	if err := s.RefreshTokenRepository.Create(
		ctx,
		userId,
		helper.HashToken(token),
		time.Now().Add(s.Cfg.EnvJwtRefreshExpiry),
	); err != nil {
		return "", err
	}

	return token, nil
}

func (s *RefreshTokenImpl) Update(ctx context.Context, oldToken string) (string, string, error) {
	newToken := helper.CreateRandomToken()

	userId, err := s.RefreshTokenRepository.Update(
		ctx,
		helper.HashToken(oldToken),
		helper.HashToken(newToken),
		time.Now().Add(s.Cfg.EnvJwtRefreshExpiry),
	)
	if err != nil {
		return "", "", helper.ErrUserInvalidOrExpiredToken
	}

	return userId, newToken, nil
}

func (s *RefreshTokenImpl) Delete(ctx context.Context, token string) error {
	if err := s.RefreshTokenRepository.Delete(ctx, helper.HashToken(token)); err != nil {
		return err
	}

	return nil
}
