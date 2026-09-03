package user

import (
	"context"
	"server/helper"
	"server/prisma/db"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, email, password, verificationToken string) (string, error)
	Delete(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	SetVerified(ctx context.Context, token string) error
	SetResetToken(ctx context.Context, email string, token string, expiry time.Time) error
	ResetPassword(ctx context.Context, token, password string) error
	UpdatePassword(ctx context.Context, id string, password string) error
}

type UserRepositoryImpl struct {
	Db *db.PrismaClient
}

func NewUserRepository(dbClient *db.PrismaClient) UserRepository {
	return &UserRepositoryImpl{Db: dbClient}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, email, password, verificationToken string) (string, error) {
	result, err := r.Db.User.
		CreateOne(
			db.User.Email.Set(email),
			db.User.Password.Set(password),
			db.User.VerificationToken.Set(verificationToken),
		).
		Exec(ctx)

	if err != nil {
		if _, ok := db.IsErrUniqueConstraint(err); ok {
			return "", helper.ErrEmailAlreadyExists
		}
		return "", err
	}

	return result.ID, nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.Db.User.
		FindUnique(db.User.ID.Equals(id)).
		Delete().
		Exec(ctx)

	if err != nil {
		if db.IsErrNotFound(err) {
			return helper.ErrUserNotFound
		}
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) FindById(ctx context.Context, id string) (User, error) {
	result, err := r.Db.User.
		FindUnique(db.User.ID.Equals(id)).
		Exec(ctx)

	if err != nil {
		if db.IsErrNotFound(err) {
			return User{}, helper.ErrUserNotFound
		}
		return User{}, err
	}

	return User{
		Id:         result.ID,
		Email:      result.Email,
		IsVerified: result.IsVerified,
	}, nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (User, error) {
	result, err := r.Db.User.
		FindUnique(db.User.Email.Equals(email)).
		Exec(ctx)

	if err != nil {
		if db.IsErrNotFound(err) {
			return User{}, helper.ErrUserNotFound
		}
		return User{}, err
	}

	verificationToken, _ := result.VerificationToken()

	return User{
		Id:                result.ID,
		Email:             result.Email,
		Password:          result.Password,
		VerificationToken: verificationToken,
		IsVerified:        result.IsVerified,
	}, nil
}

func (r *UserRepositoryImpl) SetVerified(ctx context.Context, token string) error {
	result, err := r.Db.User.
		FindMany(db.User.VerificationToken.Equals(token)).
		Update(
			db.User.VerificationToken.Set(""),
			db.User.IsVerified.Set(true),
		).
		Exec(ctx)

	if err != nil {
		return err
	}

	if result.Count == 0 {
		return helper.ErrInvalidOrExpiredToken
	}

	return nil
}

func (r *UserRepositoryImpl) SetResetToken(ctx context.Context, email, token string, expiry time.Time) error {
	_, err := r.Db.User.
		FindUnique(db.User.Email.Equals(email)).
		Update(
			db.User.ResetToken.Set(token),
			db.User.ResetTokenExpiry.Set(expiry),
		).
		Exec(ctx)

	if err != nil {
		if db.IsErrNotFound(err) {
			return helper.ErrUserNotFound
		}
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) ResetPassword(ctx context.Context, token, password string) error {
	result, err := r.Db.User.
		FindMany(
			db.User.ResetToken.Equals(token),
			db.User.ResetTokenExpiry.After(time.Now()),
		).
		Update(db.User.Password.Set(password),
			db.User.ResetToken.Set(""),
			db.User.ResetTokenExpiry.Set(time.Now()),
		).
		Exec(ctx)

	if err != nil {
		return err
	}

	if result.Count == 0 {
		return helper.ErrInvalidOrExpiredToken
	}

	return nil
}

func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, id string, password string) error {
	_, err := r.Db.User.
		FindUnique(db.User.ID.Equals(id)).
		Update(
			db.User.Password.Set(password),
		).
		Exec(ctx)

	if err != nil {
		if db.IsErrNotFound(err) {
			return helper.ErrUserNotFound
		}
		return err
	}

	return nil
}
