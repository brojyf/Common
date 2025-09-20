package middlewares

import "github.com/gin-gonic/gin"

func OneTimeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("email", "patrick.jiang@plu.edu")
		c.Set("scene", "signup")
		c.Set("jti", "")
		c.Next()
	}
}
