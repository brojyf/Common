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
