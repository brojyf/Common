package httpx

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

func IsCtxDone(ctx context.Context, err error) bool {
	if err != nil && (errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
		return true
	}
	return ctx != nil && ctx.Err() != nil
}

func ShouldSkipWrite(c *gin.Context, err error) bool {
	return c.IsAborted() || c.Writer.Written() || IsCtxDone(c.Request.Context(), err)
}
