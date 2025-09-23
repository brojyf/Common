package middlewares

import (
	"backend/internal/config"
	"backend/internal/pkg/httpx"
	"backend/internal/pkg/jwtx"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func OneTimeToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Extract Bearer
		tokenStr, err := extractBearer(c)
		if err != nil {
			httpx.WriteUnauthorized(c, "Cannot find token")
			return
		}

		// 2. Parse token
		var claims jwtx.OTTClaims
		token, err := jwt.ParseWithClaims(
			tokenStr,
			&claims,
			func(t *jwt.Token) (any, error) {
				if t.Method != jwt.SigningMethodHS256 {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return config.C.JWT.KEY, nil
			},
		)
		if err != nil || !token.Valid {
			httpx.WriteUnauthorized(c, "Token is invalid or expired")
			return
		}

		// 3. Verify Issuer and ExpireAt
		now := time.Now()
		if !verifyIssuer(&claims.RegisteredClaims, config.C.JWT.ISS) {
			httpx.WriteUnauthorized(c, "Token is invalid")
			return
		}
		if !verifyExpiresAt(&claims.RegisteredClaims, now) {
			httpx.WriteUnauthorized(c, "Token expired")
			return
		}

		c.Set("email", claims.Email)
		c.Set("scene", claims.Scene)
		c.Set("jti", claims.ID)

		c.Next()
	}
}

func verifyIssuer(rc *jwt.RegisteredClaims, expected string) bool {
	if rc == nil {
		return false
	}
	return rc.Issuer == expected
}

func verifyExpiresAt(rc *jwt.RegisteredClaims, now time.Time) bool {
	if rc == nil || rc.ExpiresAt == nil {
		return false
	}
	return rc.ExpiresAt.After(now)
}

func extractBearer(c *gin.Context) (string, error) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return "", errors.New("missing Authorization header")
	}
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", errors.New("invalid Authorization format")
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		return "", errors.New("empty bearer token")
	}
	return token, nil
}
