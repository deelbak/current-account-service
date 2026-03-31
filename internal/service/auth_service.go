package service

import (
	"context"
	"fmt"

	"current-account-service/internal/middleware"
	"current-account-service/internal/models"
	"current-account-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users repository.UserRepository
}

func NewAuthService(users repository.UserRepository) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, *models.User, error) {
	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, fmt.Errorf("token generation failed")
	}

	return token, user, nil
}
