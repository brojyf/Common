package jwt

import (
	"Backend/internal/config"
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

func parse()
