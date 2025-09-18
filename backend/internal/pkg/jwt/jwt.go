package jwt

import (
	"backend/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignOTT(email, scene string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := OTTClaims{
		Email: email,
		Scene: scene,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(config.C.JWT.KEY)
}

type OTTClaims struct {
	Email string `json:"email"`
	Scene string `json:"scene"`
	jwt.RegisteredClaims
}
