package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alesr/usrsvc/internal/users/repository"
	"github.com/alesr/usrsvc/pkg/events"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestFetch(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		// Arrange

		storedUser := &repository.User{
			ID:        uuid.New().String(),
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "jdoe",
			Password:  "password",
			Email:     "joedoe@foo.bar",
			Country:   "US",
			CreatedAt: time.Time{}.Add(1 * time.Second),
			UpdatedAt: time.Time{}.Add(2 * time.Second),
		}

		var getFuncWasCalled bool
		repo := &repoMock{
			GetFunc: func(ctx context.Context, id string) (*repository.User, error) {
				getFuncWasCalled = true
				require.Equal(t, storedUser.ID, id)
				return storedUser, nil
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualUser, err := svc.Fetch(context.TODO(), storedUser.ID)
		require.NoError(t, err)

		// Assert

		assert.True(t, getFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Equal(t, storedUser.ID, actualUser.ID)
		assert.Equal(t, storedUser.FirstName, actualUser.FirstName)
		assert.Equal(t, storedUser.LastName, actualUser.LastName)
		assert.Equal(t, storedUser.Nickname, actualUser.Nickname)
		assert.Equal(t, storedUser.Password, actualUser.Password)
		assert.Equal(t, storedUser.Email, actualUser.Email)
		assert.Equal(t, storedUser.Country, actualUser.Country)
		assert.Equal(t, storedUser.CreatedAt, actualUser.CreatedAt)
		assert.Equal(t, storedUser.UpdatedAt, actualUser.UpdatedAt)
	})

	t.Run("user not found", func(t *testing.T) {
		// Arrange

		var getFuncWasCalled bool
		repo := &repoMock{
			GetFunc: func(ctx context.Context, id string) (*repository.User, error) {
				getFuncWasCalled = true
				return nil, repository.ErrUserNotFound
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualUser, actualErr := svc.Fetch(context.TODO(), uuid.New().String())

		// Assert

		assert.Error(t, actualErr)
		assert.True(t, getFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.True(t, errors.Is(actualErr, ErrUserNotFound))
		assert.Nil(t, actualUser)
	})

	t.Run("missing id", func(t *testing.T) {
		// Arrange

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), &repoMock{}, WithPublisher(publisher))

		// Act
		actualUser, actualErr := svc.Fetch(context.TODO(), "")

		// Assert
		assert.Error(t, actualErr)
		assert.False(t, publisherWasCalled)
		assert.True(t, errors.Is(actualErr, ErrInvalidID))
		assert.Nil(t, actualUser)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), &repoMock{}, WithPublisher(publisher))

		// Act
		actualUser, actualErr := svc.Fetch(context.TODO(), "invalid-id")

		// Assert

		assert.Error(t, actualErr)
		assert.False(t, publisherWasCalled)
		assert.True(t, errors.Is(actualErr, ErrInvalidID))
		assert.Nil(t, actualUser)
	})

	t.Run("repo error", func(t *testing.T) {
		// Arrange

		var getFuncWasCalled bool
		repo := &repoMock{
			GetFunc: func(ctx context.Context, id string) (*repository.User, error) {
				getFuncWasCalled = true
				return nil, errors.New("repo error")
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualUser, actualErr := svc.Fetch(context.TODO(), uuid.New().String())

		// Assert
		assert.True(t, getFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Nil(t, actualUser)
		assert.Error(t, actualErr)
	})
}

func TestFetchAll(t *testing.T) {
	t.Parallel()

	t.Run("success without filtering", func(t *testing.T) {
		// Arrange

		id1 := uuid.New().String()
		id2 := uuid.New().String()

		var getAllFuncWasCalled bool

		repo := &repoMock{
			GetAllFunc: func(ctx context.Context, cursor string, limit int) ([]*repository.User, error) {
				getAllFuncWasCalled = true
				return []*repository.User{
					{
						ID:        id1,
						FirstName: "John",
						LastName:  "Doe",
						Nickname:  "jdoe",
						Password:  "password",
						Email:     "joedoe@foo.bar",
						Country:   "US",
						CreatedAt: time.Time{}.Add(1 * time.Second),
						UpdatedAt: time.Time{}.Add(2 * time.Second),
					},
					{
						ID:        id2,
						FirstName: "Jane",
						LastName:  "Doe",
						Nickname:  "jdoe",
						Password:  "password",
						Email:     "janedoe@foo.bar",
						Country:   "US",
						CreatedAt: time.Time{}.Add(1 * time.Second),
						UpdatedAt: time.Time{}.Add(2 * time.Second),
					},
				}, nil
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualUser, err := svc.FetchAll(context.TODO(), FilterParams{}, PaginationParams{})
		require.NoError(t, err)

		// Assert

		require.Len(t, actualUser, 2)

		assert.True(t, getAllFuncWasCalled)
		assert.False(t, publisherWasCalled)

		assert.Equal(t, id1, actualUser[0].ID)
		assert.Equal(t, "John", actualUser[0].FirstName)
		assert.Equal(t, "Doe", actualUser[0].LastName)
		assert.Equal(t, "jdoe", actualUser[0].Nickname)
		assert.Equal(t, "password", actualUser[0].Password)
		assert.Equal(t, "joedoe@foo.bar", actualUser[0].Email)
		assert.Equal(t, "US", actualUser[0].Country)
		assert.Equal(t, time.Time{}.Add(1*time.Second), actualUser[0].CreatedAt)
		assert.Equal(t, time.Time{}.Add(2*time.Second), actualUser[0].UpdatedAt)

		assert.Equal(t, id2, actualUser[1].ID)
		assert.Equal(t, "Jane", actualUser[1].FirstName)
		assert.Equal(t, "Doe", actualUser[1].LastName)
		assert.Equal(t, "jdoe", actualUser[1].Nickname)
		assert.Equal(t, "password", actualUser[1].Password)
		assert.Equal(t, "janedoe@foo.bar", actualUser[1].Email)
		assert.Equal(t, "US", actualUser[1].Country)
		assert.Equal(t, time.Time{}.Add(1*time.Second), actualUser[1].CreatedAt)
		assert.Equal(t, time.Time{}.Add(2*time.Second), actualUser[1].UpdatedAt)
	})

	t.Run("empty list", func(t *testing.T) {
		// Arrange

		var getAllFuncWasCalled bool
		repo := &repoMock{
			GetAllFunc: func(ctx context.Context, cursor string, limit int) ([]*repository.User, error) {
				getAllFuncWasCalled = true
				return []*repository.User{}, nil
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualUser, err := svc.FetchAll(context.TODO(), FilterParams{}, PaginationParams{})
		require.NoError(t, err)

		// Assert
		require.True(t, getAllFuncWasCalled)
		require.False(t, publisherWasCalled)
		assert.Len(t, actualUser, 0)
	})

	t.Run("success with filtering", func(t *testing.T) {
		// Arrange

		id1 := uuid.New().String()
		id2 := uuid.New().String()

		var getByCountryFuncWasCalled bool

		repo := &repoMock{
			GetByCountryFunc: func(ctx context.Context, country, cursor string, limit int) ([]*repository.User, error) {
				getByCountryFuncWasCalled = true
				return []*repository.User{
					{
						ID:        id1,
						FirstName: "John",
						LastName:  "Doe",
						Nickname:  "jdoe",
						Password:  "password",
						Email:     "joedoe@foo.bar",
						Country:   "US",
						CreatedAt: time.Time{}.Add(1 * time.Second),
						UpdatedAt: time.Time{}.Add(2 * time.Second),
					},
					{
						ID:        id2,
						FirstName: "Jane",
						LastName:  "Doe",
						Nickname:  "jdoe",
						Password:  "password",
						Email:     "janedoe@foo.bar",
						Country:   "US",
						CreatedAt: time.Time{}.Add(1 * time.Second),
						UpdatedAt: time.Time{}.Add(2 * time.Second),
					},
				}, nil
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act

		country := "US"
		actualUser, err := svc.FetchAll(context.TODO(), FilterParams{Country: &country}, PaginationParams{})
		require.NoError(t, err)

		// Assert

		require.Len(t, actualUser, 2)

		assert.True(t, getByCountryFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Equal(t, id1, actualUser[0].ID)
		assert.Equal(t, "John", actualUser[0].FirstName)
		assert.Equal(t, "Doe", actualUser[0].LastName)
		assert.Equal(t, "jdoe", actualUser[0].Nickname)
		assert.Equal(t, "password", actualUser[0].Password)
		assert.Equal(t, "joedoe@foo.bar", actualUser[0].Email)
		assert.Equal(t, "US", actualUser[0].Country)
		assert.Equal(t, time.Time{}.Add(1*time.Second), actualUser[0].CreatedAt)
		assert.Equal(t, time.Time{}.Add(2*time.Second), actualUser[0].UpdatedAt)

		assert.Equal(t, id2, actualUser[1].ID)
		assert.Equal(t, "Jane", actualUser[1].FirstName)
		assert.Equal(t, "Doe", actualUser[1].LastName)
		assert.Equal(t, "jdoe", actualUser[1].Nickname)
		assert.Equal(t, "password", actualUser[1].Password)
		assert.Equal(t, "janedoe@foo.bar", actualUser[1].Email)
		assert.Equal(t, "US", actualUser[1].Country)
		assert.Equal(t, time.Time{}.Add(1*time.Second), actualUser[1].CreatedAt)
		assert.Equal(t, time.Time{}.Add(2*time.Second), actualUser[1].UpdatedAt)
	})

	t.Run("missing country code", func(t *testing.T) {
		// Arrange

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), &repoMock{}, WithPublisher(publisher))

		// Act
		country := ""
		actualUser, actualErr := svc.FetchAll(
			context.TODO(),
			FilterParams{Country: &country},
			PaginationParams{},
		)

		// Assert
		assert.Error(t, actualErr)
		assert.False(t, publisherWasCalled)
		assert.True(t, errors.Is(actualErr, ErrCountryCodeInvalid))
		assert.Nil(t, actualUser)
	})

	t.Run("invalid country code", func(t *testing.T) {
		// Arrange
		svc := NewServiceDefault(zap.NewNop(), &repoMock{})

		// Act
		country := "invalid-country"
		actualUser, actualErr := svc.FetchAll(
			context.TODO(),
			FilterParams{Country: &country},
			PaginationParams{},
		)

		// Assert
		require.Error(t, actualErr)
		assert.True(t, errors.Is(actualErr, ErrCountryCodeInvalid))
		assert.Nil(t, actualUser)
	})

	t.Run("repo get all error", func(t *testing.T) {
		// Arrange

		var getAllFuncWasCalled bool

		repo := &repoMock{
			GetAllFunc: func(ctx context.Context, cursor string, limit int) ([]*repository.User, error) {
				getAllFuncWasCalled = true
				return nil, errors.New("repo error")
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualUser, actualErr := svc.FetchAll(context.TODO(),
			FilterParams{},
			PaginationParams{},
		)

		// Assert
		assert.True(t, getAllFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Error(t, actualErr)
		assert.Error(t, actualErr)
		assert.Nil(t, actualUser)
	})

	t.Run("repo get by country error", func(t *testing.T) {
		// Arrange

		var getByCountryFuncWasCalled bool

		repo := &repoMock{
			GetByCountryFunc: func(ctx context.Context, country, cursor string, limit int) ([]*repository.User, error) {
				getByCountryFuncWasCalled = true
				return nil, errors.New("repo error")
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		country := "US"
		actualUser, actualErr := svc.FetchAll(
			context.TODO(),
			FilterParams{Country: &country},
			PaginationParams{})

		// Assert

		assert.True(t, getByCountryFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Error(t, actualErr)
		assert.Error(t, actualErr)
		assert.Nil(t, actualUser)
	})
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange

		givenUser := &User{
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "jdoe",
			Password:  "password",
			Email:     "joedoe@foo.bar",
			Country:   "US",
		}

		var insertFuncWasCalled bool
		repo := &repoMock{
			InsertFunc: func(ctx context.Context, user *repository.User) error {
				insertFuncWasCalled = true

				// Assert if the values passed to the repo are as expected

				assert.NotNil(t, user.ID)
				_, err := uuid.Parse(user.ID)
				assert.NoError(t, err)

				assert.Equal(t, givenUser.FirstName, user.FirstName)
				assert.Equal(t, givenUser.LastName, user.LastName)
				assert.Equal(t, givenUser.Nickname, user.Nickname)
				assert.Equal(t, givenUser.Password, user.Password)
				assert.Equal(t, givenUser.Email, user.Email)
				assert.Equal(t, givenUser.Country, user.Country)
				assert.NotEmpty(t, user.CreatedAt)
				assert.NotEmpty(t, user.UpdatedAt)
				return nil
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualUser, err := svc.Create(context.TODO(), givenUser)
		require.NoError(t, err)

		// Assert

		assert.True(t, insertFuncWasCalled)
		assert.True(t, publisherWasCalled)

		assert.NotNil(t, actualUser.ID)
		_, err = uuid.Parse(actualUser.ID)
		assert.NoError(t, err)

		assert.Equal(t, givenUser.FirstName, actualUser.FirstName)
		assert.Equal(t, givenUser.LastName, actualUser.LastName)
		assert.Equal(t, givenUser.Nickname, actualUser.Nickname)
		assert.Equal(t, givenUser.Password, actualUser.Password)
		assert.Equal(t, givenUser.Email, actualUser.Email)
		assert.Equal(t, givenUser.Country, actualUser.Country)
		assert.NotEmpty(t, actualUser.CreatedAt)
		assert.NotEmpty(t, actualUser.UpdatedAt)
	})

	t.Run("user already exists", func(t *testing.T) {
		// Arrange

		var insertFuncWasCalled bool

		repo := &repoMock{
			InsertFunc: func(ctx context.Context, user *repository.User) error {
				insertFuncWasCalled = true
				return repository.ErrDuplicateEmail
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act

		actualUser, actualErr := svc.Create(context.TODO(), &User{})

		// Assert
		assert.True(t, insertFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Error(t, actualErr)
		assert.True(t, errors.Is(actualErr, ErrUserAlreadyExists))
		assert.Nil(t, actualUser)
	})

	t.Run("repo insert error", func(t *testing.T) {
		// Arrange

		var insertFuncWasCalled bool

		repo := &repoMock{
			InsertFunc: func(ctx context.Context, user *repository.User) error {
				insertFuncWasCalled = true
				return errors.New("repo error")
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualUser, actualErr := svc.Create(context.TODO(), &User{})

		// Assert
		assert.True(t, insertFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Error(t, actualErr)
		assert.Nil(t, actualUser)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange

		givenUser := &User{
			ID:        uuid.New().String(),
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "jdoe",
			Password:  "password",
			Email:     "joedoe@foo.bar",
			Country:   "US",
			CreatedAt: time.Time{}.Add(time.Duration(1) * time.Second),
			UpdatedAt: time.Time{}.Add(time.Duration(2) * time.Second),
		}

		var updateFuncWasCalled bool
		repo := &repoMock{
			UpdateFunc: func(ctx context.Context, user *repository.User) error {
				updateFuncWasCalled = true

				// Assert if the values passed to the repo are as expected

				assert.Equal(t, givenUser.ID, user.ID)
				assert.Equal(t, givenUser.FirstName, user.FirstName)
				assert.Equal(t, givenUser.LastName, user.LastName)
				assert.Equal(t, givenUser.Nickname, user.Nickname)
				assert.Equal(t, givenUser.Password, user.Password)
				assert.Equal(t, givenUser.Email, user.Email)
				assert.Equal(t, givenUser.Country, user.Country)
				assert.Equal(t, givenUser.CreatedAt, user.CreatedAt)
				assert.NotEmpty(t, user.UpdatedAt)
				return nil
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act

		actualUser, err := svc.Update(context.TODO(), givenUser)
		require.NoError(t, err)

		// Assert

		require.True(t, updateFuncWasCalled)
		require.True(t, publisherWasCalled)

		assert.NotNil(t, actualUser.ID)
		_, err = uuid.Parse(actualUser.ID)
		assert.NoError(t, err)

		assert.Equal(t, givenUser.FirstName, actualUser.FirstName)
		assert.Equal(t, givenUser.LastName, actualUser.LastName)
		assert.Equal(t, givenUser.Nickname, actualUser.Nickname)
		assert.Equal(t, givenUser.Password, actualUser.Password)
		assert.Equal(t, givenUser.Email, actualUser.Email)
		assert.Equal(t, givenUser.Country, actualUser.Country)
		assert.NotEmpty(t, actualUser.CreatedAt)
		assert.NotEmpty(t, actualUser.UpdatedAt)
	})

	t.Run("user not found", func(t *testing.T) {
		// Arrange

		var updateFuncWasCalled bool
		repo := &repoMock{
			UpdateFunc: func(ctx context.Context, user *repository.User) error {
				updateFuncWasCalled = true
				return repository.ErrUserNotFound
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act

		actualUser, actualErr := svc.Update(context.TODO(), &User{})

		// Assert

		assert.True(t, updateFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Error(t, actualErr)
		assert.True(t, errors.Is(actualErr, ErrUserNotFound))
		assert.Nil(t, actualUser)
	})

	t.Run("repo update error", func(t *testing.T) {
		// Arrange

		var updateFuncWasCalled bool
		repo := &repoMock{
			UpdateFunc: func(ctx context.Context, user *repository.User) error {
				updateFuncWasCalled = true
				return errors.New("repo error")
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act

		actualUser, actualErr := svc.Update(context.TODO(), &User{})

		// Assert

		assert.True(t, updateFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Error(t, actualErr)
		assert.Nil(t, actualUser)
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange

		var deleteFuncWasCalled bool
		repo := &repoMock{
			DeleteFunc: func(ctx context.Context, id string) error {
				deleteFuncWasCalled = true
				return nil
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act

		err := svc.Delete(context.TODO(), uuid.New().String())

		// Assert

		assert.True(t, deleteFuncWasCalled)
		assert.True(t, publisherWasCalled)
		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		// Arrange

		var deleteFuncWasCalled bool
		repo := &repoMock{
			DeleteFunc: func(ctx context.Context, id string) error {
				deleteFuncWasCalled = true
				return repository.ErrUserNotFound
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualErr := svc.Delete(context.TODO(), uuid.New().String())

		// Assert
		assert.True(t, deleteFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.NoError(t, actualErr)
	})

	t.Run("repo delete error", func(t *testing.T) {
		// Arrange

		var deleteFuncWasCalled bool
		repo := &repoMock{
			DeleteFunc: func(ctx context.Context, id string) error {
				deleteFuncWasCalled = true
				return errors.New("repo error")
			},
		}

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), repo, WithPublisher(publisher))

		// Act
		actualErr := svc.Delete(context.TODO(), uuid.New().String())

		// Assert
		assert.True(t, deleteFuncWasCalled)
		assert.False(t, publisherWasCalled)
		assert.Error(t, actualErr)
	})

	t.Run("invalid id", func(t *testing.T) {
		// Arrange

		var publisherWasCalled bool
		publisher := &publisherMock{
			PublishFunc: func(event events.Event, data any) error {
				publisherWasCalled = true
				return nil
			},
		}

		svc := NewServiceDefault(zap.NewNop(), &repoMock{}, WithPublisher(publisher))

		// Act
		actualErr := svc.Delete(context.TODO(), "invalid")

		// Assert
		assert.False(t, publisherWasCalled)
		assert.Error(t, actualErr)
		assert.True(t, errors.Is(actualErr, ErrInvalidID))
	})
}
