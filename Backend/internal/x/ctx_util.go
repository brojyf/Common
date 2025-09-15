package x

import (
	"context"
	"errors"
	"time"

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

func ChildWithBudget(parent context.Context, budget time.Duration) (context.Context, context.CancelFunc) {
	if dl, ok := parent.Deadline(); ok {
		remain := time.Until(dl)
		if remain <= 0 {
			c, cancel := context.WithCancel(parent)
			cancel()
			return c, func() {}
		}
		if remain < budget {
			// 留一点点余量，避免两边同时到点的竞态
			if remain > 10*time.Millisecond {
				budget = remain - 10*time.Millisecond
			} else {
				budget = remain
			}
		}
	}
	return context.WithTimeout(parent, budget)
}
