package services

import (
	"Backend/internal/config"
	"Backend/internal/repo"
	"Backend/internal/x"
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

type AuthService interface {
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
	if x.IsCtxDone(ctx, nil) {
		return ctx.Err()
	}
	c, cancel := context.WithTimeout(ctx, config.C.Timeouts.RequestCode)
	defer cancel()
	// Set throttle
	if err := s.authRepo.SetThrottle(c, email, scene); err != nil {
		if x.IsCtxDone(c, err) {
			return err
		}
		x.LogError(ctx, "AuthService.RequestCode.SetThrottle", err)
		return err
	}
	// Generate & Store Code
	code := genCode()
	if err := s.authRepo.StoreCode(c, code, email, scene); err != nil {
		if x.IsCtxDone(c, err) {
			return err
		}
		x.LogError(ctx, "AuthService.RequestCode.StoreCode", err)
		return err
	}
	// Send
	// TODO: Send via email
	log.Printf("[DEV] Verification code for %s is: %s", email, code)
	return nil
}

func genCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}
