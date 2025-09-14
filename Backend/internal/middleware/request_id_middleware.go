package middleware

import (
	"Backend/internal/pkg/request_id"
	_ "Backend/internal/pkg/request_id"

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
		ctx := request_id.With(c.Request.Context(), rid)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
