package repos

import "errors"

var (
	// ErrOTPInvalid 401
	ErrOTPInvalid = errors.New("invalid one-time password")
	// ErrOTPExpired 401
	ErrOTPExpired = errors.New("expired one-time password")

	ErrNotFound = errors.New("not found")

	// ErrEmailAlreadyExists 409
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrRateLimited 429
	ErrRateLimited = errors.New("throttle")

	// ErrRunScript 500: Error when running lua script
	ErrRunScript = errors.New("run script error")
	// ErrUnexpectedReply 500: Unexpected reply from lua
	ErrUnexpectedReply = errors.New("unexpected response")
	// ErrUnexpectedSQL 500: SQL Error
	ErrUnexpectedSQL = errors.New("unexpected error when running sql")
	// ErrUnexpectedRedis 500: Redis error
	ErrUnexpectedRedis = errors.New("unexpected error when running redis")
)
