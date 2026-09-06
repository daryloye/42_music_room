package refreshtoken

import (
	"context"
	"server/helper"
	"server/prisma/db"
	"time"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, userId, token string, expiry time.Time) error
	Update(ctx context.Context, oldToken, newToken string, expiry time.Time) (string, error)
	Delete(ctx context.Context, token string) error
}

type RefreshTokenRepositoryImpl struct {
	Db *db.PrismaClient
}

func NewRefreshTokenRepository(dbClient *db.PrismaClient) RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{Db: dbClient}
}

func (r *RefreshTokenRepositoryImpl) Create(ctx context.Context, userId, hashToken string, expiry time.Time) error {
	_, err := r.Db.RefreshToken.
		CreateOne(
			db.RefreshToken.Token.Set(hashToken),
			db.RefreshToken.Expiry.Set(expiry),
			db.RefreshToken.User.Link(
				db.User.ID.Equals(userId),
			),
		).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (r *RefreshTokenRepositoryImpl) Update(ctx context.Context, oldToken, newToken string, expiry time.Time) (string, error) {
	result, err := r.Db.RefreshToken.
		FindUnique(db.RefreshToken.Token.Equals(oldToken)).
		Exec(ctx)

	if err != nil {
		if db.IsErrNotFound(err) {
			return "", helper.ErrUserInvalidOrExpiredToken
		}
		return "", err
	}

	if result.Expiry.Before(time.Now()) {
		return "", helper.ErrUserInvalidOrExpiredToken
	}

	_, err = r.Db.RefreshToken.
		FindUnique(db.RefreshToken.Token.Equals(oldToken)).
		Update(
			db.RefreshToken.Token.Set(newToken),
			db.RefreshToken.Expiry.Set(expiry),
		).
		Exec(ctx)

	if err != nil {
		return "", err
	}

	return result.UserID, nil
}

func (r *RefreshTokenRepositoryImpl) Delete(ctx context.Context, token string) error {
	_, err := r.Db.RefreshToken.
		FindUnique(db.RefreshToken.Token.Equals(token)).
		Delete().
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
