package services

import "errors"

var (
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrConflict		  = errors.New("conflict")
	ErrTooManyRequest = errors.New("too many requests")
	ErrInternalServer = errors.New("internal server error")

	ErrCtxError = errors.New("timeout")
)
