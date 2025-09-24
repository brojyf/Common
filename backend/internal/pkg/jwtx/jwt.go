package jwtx

import (
	"backend/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignATK(uid uint64, tokenV uint, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := ATKClaims{
		UserID: uid,
		TokenV: tokenV,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.C.JWT.ISS,
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(config.C.JWT.KEY)
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

type ATKClaims struct {
	UserID uint64 `json:"email"`
	TokenV uint   `json:"token_version"`
	jwt.RegisteredClaims
}

type OTTClaims struct {
	Email string `json:"email"`
	Scene string `json:"scene"`
	jwt.RegisteredClaims
}
