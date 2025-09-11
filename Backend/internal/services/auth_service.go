package services

import (
	"Backend/internal/repo"
	"fmt"
	"log"
	"strings"
)

type AuthService interface {
	// RequestCode
	RequestCode(email, scene string) error
	CheckRequestCodeThrottle(email, scene string) error
}

type authService struct {
	authRepo repo.AuthRepo
}

func NewAuthService(authRepo repo.AuthRepo) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) CheckRequestCodeThrottle(email, scene string) error {
	if true {
		return fmt.Errorf("request code throttled")
	}
	return nil
}

func (s *authService) RequestCode(email, scene string) error {
	// Check email
	email = strings.ToLower(strings.TrimSpace(email))
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email")
	}
	// Set throttle

	// Generate & Store Code
	code := "000000"

	// Send & Return
	// TODO: Send via email
	log.Printf("[DEV] Verification code for %s is: %s", email, code)
	return nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
}
