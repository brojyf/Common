package httpx

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TryWriteJSON Respect to context, used for 2xx
func TryWriteJSON(c *gin.Context, ctx context.Context, code int, data interface{}) {
	if c.IsAborted() || c.Writer.Written() {
		return
	}
	if err := ctx.Err(); err != nil {
		WriteCtxError(c, err)
		return
	}
	c.JSON(code, data)
}

func WriteBadReq(c *gin.Context) {
	WriteJSON(c, http.StatusBadRequest, gin.H{"error": "bad request"})
}

func WriteUnauthorized(c *gin.Context) {
	WriteJSON(c, http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}

func WriteTooManyReq(c *gin.Context) {
	WriteJSON(c, http.StatusTooManyRequests, gin.H{"error": "too many request"})
}

func WriteInternal(c *gin.Context) {
	WriteJSON(c, http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

// WriteCtxError Write 499/504, default 500
func WriteCtxError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"code": "REQUEST_TIMEOUT", "error": "request timed out"})
	case errors.Is(err, context.Canceled):
		c.AbortWithStatusJSON(499, gin.H{"code": "REQUEST_CANCELED", "error": "client canceled request"})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "error": "internal server error"})
	}
}

// WriteJSON Force to write
func WriteJSON(c *gin.Context, code int, data any) {
	if c.IsAborted() || c.Writer.Written() {
		return
	}
	c.AbortWithStatusJSON(code, data)
}
