package services

import (
	"Backend/internal/repo"
	"context"
	"log"
)

type AuthService interface {
	// RequestCode
	RequestCode(ctx context.Context, email, scene string) error
	CheckRequestCodeThrottle(ctx context.Context, email, scene string) error
}

type authService struct {
	authRepo repo.AuthRepo
}

func NewAuthService(authRepo repo.AuthRepo) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) CheckRequestCodeThrottle(ctx context.Context, email, scene string) error {
	return nil
}

func (s *authService) RequestCode(ctx context.Context, email, scene string) error {
	// Set throttle

	// Generate & Store Code
	code := "000000"

	// Send & Return
	// TODO: Send via email
	log.Printf("[DEV] Verification code for %s is: %s", email, code)
	return nil
}
