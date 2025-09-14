package errors

import "errors"

// Common errors
var (
	ErrFlowNotFound   = errors.New("flow not found")
	ErrInvalidState   = errors.New("invalid state")
	ErrInvalidOption  = errors.New("invalid option")
	ErrInvalidData    = errors.New("invalid data")
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidWebhook = errors.New("invalid webhook payload")
	ErrMissingToken   = errors.New("missing verification token")
	ErrInvalidToken   = errors.New("invalid verification token")
)
