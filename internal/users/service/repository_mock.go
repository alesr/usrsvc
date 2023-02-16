package service

import (
	"context"

	"github.com/alesr/usrsvc/internal/users/repository"
)

var _ repo = (*repoMock)(nil)

// Mock is a mock implementation of the repository interface.
type repoMock struct {
	GetFunc                 func(ctx context.Context, id string) (*repository.User, error)
	GetAllFunc              func(ctx context.Context, cursor string, limit int) ([]*repository.User, error)
	GetByCountryFunc        func(ctx context.Context, country string, cursor string, limit int) ([]*repository.User, error)
	InsertFunc              func(ctx context.Context, user *repository.User) error
	UpdateFunc              func(ctx context.Context, user *repository.User) error
	DeleteFunc              func(ctx context.Context, id string) error
	CheckDatabaseHealthFunc func(ctx context.Context) error
}

func (r *repoMock) Get(ctx context.Context, id string) (*repository.User, error) {
	return r.GetFunc(ctx, id)
}

func (r *repoMock) GetAll(ctx context.Context, cursor string, limit int) ([]*repository.User, error) {
	return r.GetAllFunc(ctx, cursor, limit)
}

func (r *repoMock) GetByCountry(ctx context.Context, country string, cursor string, limit int) ([]*repository.User, error) {
	return r.GetByCountryFunc(ctx, country, cursor, limit)
}

func (r *repoMock) Insert(ctx context.Context, user *repository.User) error {
	return r.InsertFunc(ctx, user)
}

func (r *repoMock) Update(ctx context.Context, user *repository.User) error {
	return r.UpdateFunc(ctx, user)
}

func (r *repoMock) Delete(ctx context.Context, id string) error {
	return r.DeleteFunc(ctx, id)
}

func (r *repoMock) CheckDatabaseHealth(ctx context.Context) error {
	return r.CheckDatabaseHealthFunc(ctx)
}
