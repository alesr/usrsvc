package repository

import "errors"

var (
	// Enumerate all the errors that can be returned by the repository.

	ErrDuplicateEmail error = errors.New("user already exists with given email")
	ErrUserNotFound   error = errors.New("user not found")
)
