package app

import (
	"errors"
	"fmt"

	"github.com/alesr/usrsvc/internal/users/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCountryCodeInvalid  error = status.Errorf(codes.InvalidArgument, "invalid country")
	ErrCountryCodeRequired error = status.Errorf(codes.Internal, "country is required")
	ErrEmailFormat         error = status.Errorf(codes.Internal, "email is invalid")
	ErrEmailRequired       error = status.Errorf(codes.Internal, "email is required")
	ErrIDFormat            error = status.Errorf(codes.Internal, "id is invalid")
	ErrIDRequired          error = status.Errorf(codes.Internal, "id is required")
	ErrInternal            error = status.Errorf(codes.Internal, "internal error")
	ErrNameFormat          error = status.Errorf(codes.Internal, "name must only contain letters and spaces")
	ErrNameLength          error = status.Errorf(codes.Internal, fmt.Sprintf("name must be between %d and %d characters", minNameLength, maxNameLength))
	ErrNameRequired        error = status.Errorf(codes.Internal, "name is required")
	ErrPageTokenInvalid    error = status.Errorf(codes.InvalidArgument, "invalid page token")
	ErrPasswordFormat      error = status.Errorf(codes.Internal, "password must contain at least one letter, one number and one special character")
	ErrPasswordLength      error = status.Errorf(codes.Internal, fmt.Sprintf("password must be between %d and %d characters", minPasswordLength, maxPasswordLength))
	ErrPasswordRequired    error = status.Errorf(codes.Internal, "password is required")
	ErrUserAlreadyExists   error = status.Errorf(codes.AlreadyExists, "user already exists")
	ErrUserNotFound        error = status.Errorf(codes.NotFound, "user not found")
)

func convertServiceError(svcErr error) error {
	switch {
	case errors.Is(svcErr, service.ErrCountryCodeInvalid):
		return ErrCountryCodeInvalid
	case errors.Is(svcErr, service.ErrUserNotFound):
		return ErrUserNotFound
	case errors.Is(svcErr, service.ErrUserAlreadyExists):
		return ErrUserAlreadyExists
	default:
		return ErrInternal
	}
}
