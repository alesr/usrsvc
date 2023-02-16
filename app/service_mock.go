package app

import (
	"context"

	"github.com/alesr/usrsvc/internal/users/service"
)

var _ userService = (*serviceMock)(nil)

type serviceMock struct {
	FetchFunc              func(ctx context.Context, id string) (*service.User, error)
	FetchAllFunc           func(ctx context.Context, filter service.FilterParams, pag service.PaginationParams) ([]*service.User, error)
	CreateFunc             func(ctx context.Context, user *service.User) (*service.User, error)
	UpdateFunc             func(ctx context.Context, user *service.User) (*service.User, error)
	DeleteFunc             func(ctx context.Context, id string) error
	CheckServiceHealthFunc func(ctx context.Context) error
}

func (s *serviceMock) Fetch(ctx context.Context, id string) (*service.User, error) {
	return s.FetchFunc(ctx, id)
}

func (s *serviceMock) FetchAll(ctx context.Context, filter service.FilterParams, pag service.PaginationParams) ([]*service.User, error) {
	return s.FetchAllFunc(ctx, filter, pag)
}

func (s *serviceMock) Create(ctx context.Context, user *service.User) (*service.User, error) {
	return s.CreateFunc(ctx, user)
}

func (s *serviceMock) Update(ctx context.Context, user *service.User) (*service.User, error) {
	return s.UpdateFunc(ctx, user)
}

func (s *serviceMock) Delete(ctx context.Context, id string) error {
	return s.DeleteFunc(ctx, id)
}

func (s *serviceMock) CheckServiceHealth(ctx context.Context) error {
	return s.CheckServiceHealthFunc(ctx)
}
