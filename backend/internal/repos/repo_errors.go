package repos

import "errors"

var (
	// 429
	ErrRateLimited = errors.New("throttle")

	//401
	ErrOTPInvalid = errors.New("invalid one-time password")
	ErrOTPExpired = errors.New("expired one-time password")

	// 500
	ErrRunScript       = errors.New("run script error")
	ErrUnexpectedReply = errors.New("unexpected response")
)
