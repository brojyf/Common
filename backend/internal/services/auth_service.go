package services

import (
	"backend/internal/config"
	"backend/internal/pkg/ctx_util"
	"backend/internal/pkg/jwtx"
	"backend/internal/pkg/logx"
	"backend/internal/repos"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, ip, email, password, deviceID string) (AuthResponse, error)
	CreateAccount(ctx context.Context, email, scene, jti, pwd, deviceID string) (AuthResponse, error)
	VerifyCodeAndGenToken(ctx context.Context, email, scene, codeID, code string) (string, error)
	RequestCode(ctx context.Context, email, scene string) (string, error)
}
type AuthResponse struct {
	ATK       string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
	RTK       string `json:"refresh_token"`
	UserID    uint64 `json:"user_id"`
}

func (s *authService) Login(ctx context.Context, ip, email, password, deviceID string) (AuthResponse, error) {

	// 0. Create sub context
	cctx, cancel := context.WithTimeout(ctx, config.C.Timeouts.Login)
	defer cancel()

	// 1. Check input
	if !isValidEmail(email) || !isValidPassword(password) || !isUUID(deviceID) {
		return AuthResponse{}, ErrBadRequest
	}

	// 2. Check login throttle
	rl, err := s.repo.CheckLoginThrottle(cctx, ip, email)
	if err != nil {
		if ctx_util.IsCtxDone(cctx, err) {
			return AuthResponse{}, ErrCtxError
		}
		logx.LogError(ctx, "AuthSvc.CheckLoginThrottle", err)
		return AuthResponse{}, ErrInternalServer
	}
	if rl {
		return AuthResponse{}, ErrTooManyRequest
	}

	// 3. Check password
	user, err := s.repo.GetUser(ctx, email)
	if err != nil {
		if ctx_util.IsCtxDone(cctx, err) {
			return AuthResponse{}, ErrCtxError
		}
		switch {
		default:
			logx.LogError(ctx, "AuthSvc.CreateAccount.Login", err)
			return AuthResponse{}, ErrInternalServer
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PwdHash), []byte(password))
	if err != nil {
		return AuthResponse{}, ErrUnauthorized
	}

	// 4. Update redis cnt/lock
	s.repo.ClearLoginCntLock(ip, email)

	// 5. Sign tokens
	ttl := config.C.JWT.ATK
	ttlDur := time.Duration(ttl) * time.Second
	var atk string
	for i := 0; i < 3; i++ {
		atk, err = jwtx.SignATK(user.UserID, user.TokenV, ttlDur)
		if err == nil {
			break
		}
	}
	if err != nil {
		logx.LogError(ctx, "AuthSvc.CreateAccount.SignATK", err)
		return AuthResponse{}, ErrInternalServer
	}
	rtk, rtkHash, err := generateRTK()
	if err != nil {
		logx.LogError(ctx, "AuthSvc.CreateAccount.GenerateRTK", err)
		return AuthResponse{}, ErrInternalServer
	}

	// 6. Store device id and session
	u, err := uuid.Parse(deviceID)
	if err != nil {
		logx.LogError(ctx, "AuthSvc.CreateAccount.ParseUUID", err)
		return AuthResponse{}, ErrInternalServer
	}
	didByte := u[:]
	exp := time.Now().Add(config.C.JWT.RTK)
	if err := s.repo.StoreDIDAndSession(cctx, user.UserID, didByte, nil, rtkHash, user.TokenV, exp); err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return AuthResponse{}, ErrCtxError
		}
	}

	return AuthResponse{
		ATK:       atk,
		TokenType: "Bearer",
		ExpiresIn: 3600,
		RTK:       rtk,
		UserID:    0,
	}, nil
}

func (s *authService) CreateAccount(ctx context.Context, email, scene, jti, pwd, deviceID string) (AuthResponse, error) {

	// 0. Create sub context
	cctx, cancel := context.WithTimeout(ctx, config.C.Timeouts.CreateAccount)
	defer cancel()

	// 1. Check input
	if !isValidPassword(pwd) || scene != "signup" {
		return AuthResponse{}, ErrBadRequest
	}

	// 2. Check conflict
	exist, err := s.repo.CheckEmailExists(cctx, email)
	if err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return AuthResponse{}, ErrCtxError
		}
		logx.LogError(ctx, "AuthSvc.CreateAccount.CheckEmailExists", err)
		return AuthResponse{}, ErrInternalServer
	}
	if exist {
		return AuthResponse{}, ErrConflict
	}

	// 3. Check and mark jti
	newTTL := int(config.C.JWT.OTT.Seconds())
	if err := s.repo.ConsumeOTTJTI(cctx, email, scene, jti, newTTL); err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return AuthResponse{}, ErrCtxError
		}
		if errors.Is(err, repos.ErrOTPInvalid) {
			return AuthResponse{}, ErrUnauthorized
		}
		logx.LogError(ctx, "AuthSvc.CreateAccount.ConsumeOTTJTI", err)
		return AuthResponse{}, ErrInternalServer
	}

	// 4. Create user
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		logx.LogError(ctx, "AuthSvc.CreateAccount.Hash", err)
		return AuthResponse{}, ErrInternalServer
	}
	pwdHashStr := string(pwdHash)
	uid, tkv, err := s.repo.CreateUser(cctx, email, pwdHashStr)
	if err != nil {
		// Rollback
		s.repo.UndoOTTMark(cctx, email, scene, jti, newTTL)
		if ctx_util.IsCtxDone(ctx, err) {
			return AuthResponse{}, ErrCtxError
		}
		if errors.Is(err, repos.ErrEmailAlreadyExists) {
			return AuthResponse{}, ErrConflict
		}
		logx.LogError(ctx, "AuthSvc.CreateAccount.CreateUser", err)
		return AuthResponse{}, ErrInternalServer
	}

	// 5. Sign token
	ttl := config.C.JWT.ATK
	ttlDur := time.Duration(ttl) * time.Second
	var atk string
	for i := 0; i < 3; i++ {
		atk, err = jwtx.SignATK(uid, tkv, ttlDur)
		if err == nil {
			break
		}
	}
	if err != nil {
		logx.LogError(ctx, "AuthSvc.CreateAccount.SignATK", err)
		return AuthResponse{}, ErrInternalServer
	}
	rtk, rtkHash, err := generateRTK()
	if err != nil {
		logx.LogError(ctx, "AuthSvc.CreateAccount.GenerateRTK", err)
		return AuthResponse{}, ErrInternalServer
	}

	// 6. Store device id and session
	u, err := uuid.Parse(deviceID)
	if err != nil {
		logx.LogError(ctx, "AuthSvc.CreateAccount.ParseUUID", err)
		return AuthResponse{}, ErrInternalServer
	}
	didByte := u[:]
	exp := time.Now().Add(config.C.JWT.RTK)
	if err := s.repo.StoreDIDAndSession(cctx, uid, didByte, nil, rtkHash, tkv, exp); err != nil {
		if ctx_util.IsCtxDone(ctx, err) {
			return AuthResponse{}, ErrCtxError
		}
	}

	return AuthResponse{
		atk,
		"Bearer",
		config.C.JWT.ATK,
		rtk,
		uid,
	}, nil
}

func (s *authService) VerifyCodeAndGenToken(ctx context.Context, email, scene, codeID, code string) (string, error) {

	// 0. Create sub context
	cctx, cancel := context.WithTimeout(ctx, config.C.Timeouts.VerifyCode)
	defer cancel()

	// 1. Check email & scene
	email = strings.ToLower(strings.TrimSpace(email))
	scene = strings.TrimSpace(scene)
	if !isValidEmail(email) || !isValidScene(scene) || len(code) != 6 || !isUUID(codeID) {
		return "", ErrBadRequest
	}

	// 32. Generate JTI
	jti := uuid.NewString()

	// 3. Sign token
	ttl := config.C.JWT.OTT
	token, err := jwtx.SignOTT(email, scene, jti, ttl)
	if err != nil {
		logx.LogError(ctx, "AuthSvc.VerifyCodeAndGenToken.SignOTT", err)
		return "", ErrInternalServer
	}

	// 4. Call repo: Check throttle -> Match code -> Consume code -> Set jti unused
	verifyLimit := config.C.RedisTTL.VerifyWindowLimit
	window := config.C.RedisTTL.VerifyWindow
	jtiUsedTTL := int(config.C.JWT.OTT.Seconds())
	if s.repo.ThrottleMatchAndConsumeCode(cctx, email, scene, codeID, code, jti, verifyLimit, window, jtiUsedTTL) != nil {
		if ctx_util.IsCtxDone(cctx, err) {
			return "", ErrCtxError
		}
		switch {
		case errors.Is(err, repos.ErrOTPInvalid) || errors.Is(err, repos.ErrOTPExpired):
			return "", ErrUnauthorized
		case errors.Is(err, repos.ErrRateLimited):
			return "", ErrTooManyRequest
		default:
			logx.LogError(ctx, "AuthSvc.VerifyCodeAndGenToken.ThrottleMatchAndConsumeCode", err)
			return "", ErrInternalServer
		}
	}

	return token, nil
}

func (s *authService) RequestCode(ctx context.Context, email, scene string) (string, error) {

	// 0. Create sub context
	cctx, cancel := context.WithTimeout(ctx, config.C.Timeouts.RequestCode)
	defer cancel()

	// 1. Check email & scene
	email = strings.ToLower(strings.TrimSpace(email))
	scene = strings.TrimSpace(scene)
	if !isValidEmail(email) || !isValidScene(scene) {
		return "", ErrBadRequest
	}

	// 2. Generate code & code id
	code, err := generateCode()
	if err != nil {
		logx.LogError(ctx, "AuthSvc.GenerateCode", err)
		return "", ErrInternalServer
	}
	codeID := uuid.NewString()

	// 3. Call repo: Check throttle -> Set throttle -> Store code
	otpTTL := config.C.RedisTTL.OTP
	throttleTTL := config.C.RedisTTL.OTPThrottle
	throttled, err := s.repo.StoreOTPAndThrottle(cctx, email, scene, codeID, code, otpTTL, throttleTTL)
	if err != nil {
		if ctx_util.IsCtxDone(cctx, err) {
			return "", ErrCtxError
		}
		logx.LogError(ctx, "AuthSvc.RequestCode.StoreOTPAndThrottle", err)
		return "", ErrInternalServer
	}
	if throttled {
		return "", ErrTooManyRequest
	}

	// 4. Send code
	log.Printf("[DEV] verification code for %s: %s", email, code)
	return codeID, nil
}

func generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}
func isValidPassword(pwd string) bool {
	if len(pwd) < 8 || len(pwd) > 20 {
		return false
	}
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	for _, c := range pwd {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasLower && hasUpper && hasDigit && hasSpecial
}
func isValidEmail(email string) bool {
	var (
		localPartRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+$`)
		domainRegex    = regexp.MustCompile(`^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	)

	// net/email
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	// 0. Get addr
	address := addr.Address
	parts := strings.Split(address, "@")
	if len(parts) != 2 {
		return false
	}
	localPart, domain := parts[0], parts[1]
	// 1. Simple characters
	if !localPartRegex.MatchString(localPart) {
		return false
	}
	// 2. No IP
	if !domainRegex.MatchString(domain) {
		return false
	}
	if net.ParseIP(domain) != nil {
		return false
	}
	// 3. Length
	if len(domain) > 255 {
		return false
	}

	return true
}
func isValidScene(scene string) bool {
	return scene == "signup" || scene == "reset_password"
}
func isValidUsername(s string) bool {
	var usernameRe = regexp.MustCompile(`^[A-Za-z0-9 ]+$`)
	if len(s) == 0 {
		return false
	}
	if utf8.RuneCountInString(s) > 20 {
		return false
	}
	return usernameRe.MatchString(s)
}
func isUUID(s string) bool {
	u, err := uuid.Parse(s)
	if err != nil {
		return false
	}
	return u.Version() == 4 && u.Variant() == uuid.RFC4122
}
func generateRTK() (rtk string, rtkHash []byte, err error) {
	// 1) 生成 32 字节随机数
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", nil, err
	}

	// 2) 明文 token 给客户端（base64url 无填充，安全）
	rtk = base64.RawURLEncoding.EncodeToString(raw)

	// 3) SHA-256 摘要，落库 VARBINARY(32)
	sum := sha256.Sum256([]byte(rtk))
	rtkHash = sum[:]

	return rtk, rtkHash, nil
}

type authService struct {
	repo repos.AuthRepo
}

func NewAuthService(repo repos.AuthRepo) AuthService {
	return &authService{repo: repo}
}
