package jwt

import (
	"Backend/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret []byte
}

type AccessTokenClaims struct {
	UserID       uint64 `json:"sub"`
	DeviceID     string `json:"device_id"`
	TokenVersion uint   `json:"token_version"`
	jwt.RegisteredClaims
}

func SignATK(userID uint64, deviceID string, tokenVersion uint) (string, error) {
	now := time.Now()
	claims := AccessTokenClaims{
		UserID:       userID,
		DeviceID:     deviceID,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(config.C.JWT.ATK) * time.Second)),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(config.C.JWT.KEY)
}

func ParseATK(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&AccessTokenClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return config.C.JWT.KEY, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

type OTPClaims struct {
	Email string `json:"email"`
	Scene string `json:"scene"`
	jwt.RegisteredClaims
}

func SignOTP(email, scene, jti string) (string, error) {
	now := time.Now()
	claims := OTPClaims{
		Email: email,
		Scene: scene,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(now.Add(config.C.JWT.OTP)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(config.C.JWT.KEY)
}

func ParseOTP(tokenStr string) (*OTPClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&OTPClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return config.C.JWT.KEY, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*OTPClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
