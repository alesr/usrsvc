package service

import "errors"

var (
	// Enumerate all the errors that can be returned by the service.

	ErrInvalidID          error = errors.New("invalid id")
	ErrUserNotFound       error = errors.New("user not found")
	ErrCountryCodeInvalid error = errors.New("invalid country code")
	ErrUserAlreadyExists  error = errors.New("user already exists")
)
