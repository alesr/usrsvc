package service

import "errors"

var (
	// Enumerate all the errors that can be returned by the service.

	ErrCountryCodeInvalid error = errors.New("invalid country code")
	ErrInvalidID          error = errors.New("invalid id")
	ErrUserAlreadyExists  error = errors.New("user already exists")
	ErrUserNotFound       error = errors.New("user not found")
)
