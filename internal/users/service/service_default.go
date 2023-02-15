package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alesr/usrsvc/internal/users/repository"
	"github.com/alesr/usrsvc/pkg/events"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbTimeout time.Duration = 5 * time.Second

	// Enumerate user entity change events.
	// These events are used to publish messages to a message broker.

)

// repo is the interface that provides the repository methods
type repo interface {
	Get(ctx context.Context, id string) (*repository.User, error)
	GetAll(ctx context.Context, cursor string, limit int) ([]*repository.User, error)
	GetByCountry(ctx context.Context, country string, cursor string, limit int) ([]*repository.User, error)
	Insert(ctx context.Context, user *repository.User) error
	Update(ctx context.Context, user *repository.User) error
	Delete(ctx context.Context, id string) error
}

// ServiceDefault is the default implementation of the service interface.
type ServiceDefault struct {
	logger    *zap.Logger
	repo      repo
	publisher Publisher
}

// Publisher is the interface that provides the publish method.
type Publisher interface {
	Publish(event events.Event, data any) error
}

// Option is a function that configures the service.
type Option func(*ServiceDefault)

// WithPublisher configures the service with a publisher.
func WithPublisher(publisher Publisher) Option {
	return func(s *ServiceDefault) {
		s.publisher = publisher
	}
}

// NewServiceDefault creates a new service.
func NewServiceDefault(logger *zap.Logger, repo repo, opts ...Option) *ServiceDefault {
	s := &ServiceDefault{
		logger: logger,
		repo:   repo,
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Get returns a user by id.
func (s *ServiceDefault) Fetch(ctx context.Context, id string) (*User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("could not validate id '%s': %w", id, ErrInvalidID)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	user, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("could not fetch user with id '%s': %w", id, ErrUserNotFound)
		}

		return nil, fmt.Errorf("could not fetch user with id '%s': %w", id, err)
	}
	return newUserDomainFromStore(user), nil
}

// FetchAll returns all users or users filtered by country.
func (s *ServiceDefault) FetchAll(ctx context.Context, filter FilterParams, pag PaginationParams) ([]*User, error) {
	filter.normalize()

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	var (
		users []*repository.User
		err   error
	)

	switch {
	case filter.Country != nil:
		if err := filter.validate(); err != nil {
			return nil, fmt.Errorf("could not validate fetch all filter: %w", err)
		}

		s.logger.Debug("fetching users by country", zap.String("country", *filter.Country))

		users, err = s.repo.GetByCountry(ctx, *filter.Country, pag.Cursor, pag.Limit)
		if err != nil {
			return nil, fmt.Errorf("could not fetch users by country: %w", err)
		}
	default:
		s.logger.Debug("fetching all users")

		users, err = s.repo.GetAll(ctx, pag.Cursor, pag.Limit)
		if err != nil {
			return nil, fmt.Errorf("could not fetch users: %w", err)
		}
	}

	var usersDomain []*User
	for _, user := range users {
		usersDomain = append(usersDomain, newUserDomainFromStore(user))
	}
	return usersDomain, nil
}

// Create creates a new user.
// NOTE: I left the input validation only in the transport layer, but it could be done here too.
func (s *ServiceDefault) Create(ctx context.Context, user *User) (*User, error) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %s", err)
	}

	// Replace the password with the hash.
	user.Password = string(hash)

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if err := s.repo.Insert(ctx, newUserStoreFromDomain(user)); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, fmt.Errorf("could not insert user: %w", ErrUserAlreadyExists)
		}
		return nil, fmt.Errorf("could not insert user: %w", err)
	}

	if s.publisher != nil {
		// Just keeping it simple. The most important thing is to not publish the user's password.
		s.publisher.Publish(events.UserCreated, user.ID)
	}
	return user, nil
}

// Update updates an existing user.
// NOTE: I left the input validation only in the transport layer, but it could be done here too.
func (s *ServiceDefault) Update(ctx context.Context, user *User) (*User, error) {
	user.UpdatedAt = time.Now()

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %s", err)
	}

	// Replace the password with the hash.
	user.Password = string(hash)

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if err := s.repo.Update(ctx, newUserStoreFromDomain(user)); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, fmt.Errorf("could not update user: %w", ErrUserNotFound)
		}
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, fmt.Errorf("could not update user: %w", ErrUserAlreadyExists)
		}
		return nil, fmt.Errorf("could not update user: %w", err)
	}

	if s.publisher != nil {
		s.publisher.Publish(events.UserUpdated, user.ID)
	}
	return user, nil
}

// Delete deletes an existing user.
func (s *ServiceDefault) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("could not validate id '%s': %w", id, ErrInvalidID)
	}

	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.logger.Info("could not delete user non existing user", zap.String("id", id), zap.Error(err))
			return nil
		}
		return fmt.Errorf("could not delete user with id '%s': %w", id, err)
	}

	if s.publisher != nil {
		s.publisher.Publish(events.UserDeleted, id)
	}
	return nil
}
