package ctx_util

import (
	"context"
	"errors"
)

func IsCtxDone(ctx context.Context, err error) bool {
	// Clear context error
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}
	// Check context
	if ctx != nil {
		if e := ctx.Err(); errors.Is(e, context.DeadlineExceeded) || errors.Is(e, context.Canceled) {
			return true
		}
	}
	return false
}
