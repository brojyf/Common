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
	cctx, cancel := x.ChildWithBudget(ctx, config.C.Timeouts.RequestCode)
	defer cancel()

	has, err := s.authRepo.CheckThrottle(cctx, email, scene)
	if err != nil {
		if x.IsCtxDone(cctx, err) {
			return err
		}
		x.LogError(ctx, "AuthService.CheckRequestCodeThrottle", err)
		return err
	}
	if has {
		return fmt.Errorf("throttle")
	}
	return nil
}

func (s *authService) RequestCode(ctx context.Context, email, scene string) error {
	cctx, cancel := x.ChildWithBudget(ctx, config.C.Timeouts.RequestCode)
	defer cancel()
	// Set throttle
	if err := s.authRepo.SetThrottle(cctx, email, scene); err != nil {
		if x.IsCtxDone(cctx, err) {
			return err
		}
		x.LogError(ctx, "AuthService.RequestCode.SetThrottle", err)
		return err
	}
	// Generate & Store Code
	code := genCode()
	if err := s.authRepo.StoreCode(cctx, code, email, scene); err != nil {
		if x.IsCtxDone(cctx, err) {
			return err
		}
		x.LogError(ctx, "AuthService.RequestCode.StoreCode", err)
		return err
	}
	// Send Async
	// TODO: Send via email
	log.Printf("[DEV] Verification code for %s is: %s", email, code)
	return nil
}

func genCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}
