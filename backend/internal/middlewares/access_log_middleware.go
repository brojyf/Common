package middlewares

import "github.com/gin-gonic/gin"

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
