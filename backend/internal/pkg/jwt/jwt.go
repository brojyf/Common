package jwt

import (
	"backend/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignATK() (string, error) {
	return "", nil
}

func SignOTT(email, scene, jti string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := OTTClaims{
		Email: email,
		Scene: scene,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.C.JWT.ISS,
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			ID:        jti,
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(config.C.JWT.KEY)
}

type ATKClaims struct{}

type OTTClaims struct {
	Email string `json:"email"`
	Scene string `json:"scene"`
	jwt.RegisteredClaims
}
