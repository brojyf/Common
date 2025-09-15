package middleware

import (
	"Backend/internal/x/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func OTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) 从 Header 里取 Token
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		// 2) 解析 Token
		claims, err := jwt.ParseOTP(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// 3) 存 email 和 scene 到 gin.Context
		c.Set("email", claims.Email)
		c.Set("scene", claims.Scene)
		c.Set("jti", claims.ID)

		// 放行
		c.Next()
	}
}

func ATK() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
