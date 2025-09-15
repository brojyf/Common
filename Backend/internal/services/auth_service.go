package services

import (
	"Backend/internal/config"
	"Backend/internal/repo"
	"Backend/internal/x"
	"Backend/internal/x/jwt"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignARTK(ctx context.Context, uid uint64, deviceID string) (string, string, error)
	StoreDeviceID(ctx context.Context, deviceID string, uid uint64) error
	CreateAccount(ctx context.Context, email, password string) (uint64, error)
	SignOTP(ctx context.Context, email, scene, jti string) (string, error)
	VerifyCode(ctx context.Context, email, scene, code, jti string) (bool, error)
	RequestCode(ctx context.Context, email, scene, jti string) error
	CheckRequestCodeThrottle(ctx context.Context, email, scene string) error
}

type authService struct {
	authRepo repo.AuthRepo
}

func NewAuthService(authRepo repo.AuthRepo) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) SignARTK(ctx context.Context, uid uint64, deviceID string) (string, string, error) {
	// Sub context
	cctx, cancel := x.ChildWithBudget(ctx, config.C.Timeouts.CreateAccount)
	defer cancel()
	// Get token version
	tkVersion, err := s.authRepo.GetTKVersion(cctx, uid)
	if err != nil {
		if x.IsCtxDone(cctx, err) {
			return "", "", err
		}
		x.LogError(ctx, "AuthService.SignARTK.GetTKVersion", err)
		return "", "", err
	}
	// Sign ATK
	did, err := toUUIDBytes16(deviceID)
	if err != nil {
		if x.IsCtxDone(cctx, err) {
			return "", "", err
		}
		x.LogError(ctx, "AuthService.SignARTK.ToUUIDBytes16", err)
		return "", "", err
	}
	atk, err := jwt.SignATK(uid, deviceID, tkVersion)
	if err != nil {
		if x.IsCtxDone(cctx, err) {
			return "", "", err
		}
		x.LogError(ctx, "AuthService.SignARTK.SignATK", err)
		return "", "", err
	}
	// Sign RTK & Store
	rtk, rtkHashed, err := GenerateRTK32()
	if err != nil {
		if x.IsCtxDone(cctx, err) {
			return "", "", err
		}
		x.LogError(ctx, "AuthService.SignARTK.GenerateRTK32", err)
		return "", "", err
	}
	if err := s.authRepo.StoreRTK(cctx, uid, did, rtkHashed); err != nil {
		if x.IsCtxDone(cctx, err) {
			return "", "", err
		}
		x.LogError(ctx, "AuthService.SignARTK.StoreRTK", err)
		return "", "", err
	}

	return atk, rtk, nil
}

func (s *authService) StoreDeviceID(ctx context.Context, deviceID string, uid uint64) error {
	cctx, cancel := x.ChildWithBudget(ctx, config.C.Timeouts.CreateAccount)
	defer cancel()

	did, err := toUUIDBytes16(deviceID)
	if err != nil {
		if x.IsCtxDone(cctx, err) {
			return err
		}
		x.LogError(ctx, "AuthService.StoreDeviceID.Convert2[]byte", err)
		return err
	}
	if err := s.authRepo.StoreDeviceID(cctx, uid, did); err != nil {
		if x.IsCtxDone(cctx, err) {
			return err
		}
		x.LogError(ctx, "AuthService.StoreDeviceID", err)
		return err
	}

	return nil
}

func (s *authService) CreateAccount(ctx context.Context, email, password string) (uint64, error) {
	cctx, cancel := x.ChildWithBudget(ctx, config.C.Timeouts.CreateAccount)
	defer cancel()
	// Hash pwd
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if x.IsCtxDone(cctx, err) {
			return 0, err
		}
		x.LogError(ctx, "AuthService.CreateAccount.Hash", err)
		return 0, err
	}
	// DB
	uid, err := s.authRepo.StoreNewUser(cctx, email, string(hashedBytes))
	if err != nil {
		if x.IsCtxDone(cctx, err) {
			return 0, err
		}
		x.LogError(ctx, "AuthService.CreateAccount.StoreNewUser", err)
		return 0, err
	}

	return uid, nil
}

func (s *authService) SignOTP(ctx context.Context, email, scene, jti string) (string, error) {
	token, err := jwt.SignOTP(email, scene, jti)
	if err != nil {
		x.LogError(ctx, "AuthService.SIgnOTP", err)
		return "", err
	}
	return token, nil
}

func (s *authService) VerifyCode(ctx context.Context, email, scene, code, codeID string) (bool, error) {
	cctx, cancel := x.ChildWithBudget(ctx, config.C.Timeouts.VerifyCode)
	defer cancel()

	pass, err := s.authRepo.MatchAndConsumeOTP(cctx, email, scene, code, codeID)
	if err != nil {
		if x.IsCtxDone(cctx, err) {
			return false, err
		}
		x.LogError(ctx, "AuthService.VerifyCode", err)
		return false, err
	}
	if !pass {
		return false, errors.New("invalid or expired token")
	}
	return true, nil
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

func (s *authService) RequestCode(ctx context.Context, email, scene, id string) error {
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
	if err := s.authRepo.StoreCode(cctx, code, email, scene, id); err != nil {
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

func toUUIDBytes16(s string) ([]byte, error) {
	if u, err := uuid.Parse(s); err == nil {
		b := u[:] // [16]byte -> []byte
		return b, nil
	}
	if len(s) == 32 {
		b, err := hex.DecodeString(s)
		if err == nil && len(b) == 16 {
			return b, nil
		}
	}
	return nil, errors.New("expect UUID (with dashes) or 32-hex")
}

func GenerateRTK(nBytes int) (token string, hash []byte, err error) {
	if nBytes < 32 {
		return "", nil, fmt.Errorf("nBytes too small, want >= 32")
	}

	buf := make([]byte, nBytes)
	if _, err = rand.Read(buf); err != nil {
		return "", nil, fmt.Errorf("rand.Read: %w", err)
	}

	// Base64URL（无 '+' '/' '='），适合放在Cookie/URL里
	token = base64.RawURLEncoding.EncodeToString(buf)

	h := sha256.Sum256([]byte(token))
	hash = h[:] // 32字节

	return token, hash, nil
}

// 便捷版：固定 32 字节强度
func GenerateRTK32() (string, []byte, error) {
	return GenerateRTK(32)
}
