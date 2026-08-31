package user

import (
	"context"
	"server/helper"
)

type UserService interface {
	CreateAccount(ctx context.Context, request CreateAccountRequest) (string, error)
}

type UserServiceImpl struct {
	UserRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return &UserServiceImpl{UserRepository: userRepository}
}

func (s *UserServiceImpl) CreateAccount(ctx context.Context, request CreateAccountRequest) (string, error) {
	hashPassword, err := helper.HashPassword(request.Password)
	if err != nil {
		return "", err
	}

	userId, err := s.UserRepository.Create(ctx, request.Email, hashPassword)
	if err != nil {
		return "", err
	}

	return userId, nil
}
