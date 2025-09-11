package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const headerXRequestID = "X-Request-ID"

type ctxKeyRID struct{}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {

		rid := uuid.NewString()
		c.Writer.Header().Set(headerXRequestID, rid)

		// gin.Context
		c.Set(headerXRequestID, rid)

		// context
		ctx := context.WithValue(c.Request.Context(), ctxKeyRID{}, rid)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxKeyRID{})
	rid, ok := v.(string)
	return rid, ok
}
