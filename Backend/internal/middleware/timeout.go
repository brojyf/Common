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
					panicChan <- r
				}
				close(done)
			}()
			c.Next()
		}()

		select {
		case <-ctx.Done():
			// 超时，直接返回 504
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": "request timeout",
			})
		case p := <-panicChan:
			panic(p) // 让 gin.Recovery() 处理
		case <-done:
			// 正常执行完
		}
	}
}
