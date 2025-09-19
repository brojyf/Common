package repos

import "errors"

var (
	// ErrRateLimited 429
	ErrRateLimited = errors.New("throttle")

	// ErrOTPInvalid 401
	ErrOTPInvalid = errors.New("invalid one-time password")
	// ErrOTPExpired 401
	ErrOTPExpired = errors.New("expired one-time password")

	// ErrRunScript 500: Error when running lua script
	ErrRunScript = errors.New("run script error")
	// ErrUnexpectedReply 500: Unexpected reply from lua
	ErrUnexpectedReply = errors.New("unexpected response")
)
