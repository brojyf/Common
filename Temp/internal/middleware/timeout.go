package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Create ctx with time out
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		panicChan := make(chan any, 1)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					select {
					case panicChan <- r:
					default:
					}
					return
				}
				close(done)
			}()
			c.Next()
		}()

		select {
		case <-ctx.Done():
			// 504: Timeout
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "request timeout",
			})
		case p := <-panicChan:
			panic(p)
		case <-done:
		}
	}
}
