package repos

import "errors"

var (
	ErrRateLimited      = errors.New("throttle")
	ErrOTPInvalid       = errors.New("invalid one time password")
	ErrRunScript        = errors.New("run script error")
	ErrUnexpectedReply  = errors.New("unexpected response")
	ErrRepoUnauthorized = errors.New("repo unauthorized")
)
