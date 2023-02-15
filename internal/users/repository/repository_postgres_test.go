//go:build integration
// +build integration

package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	db := setupDBHelper(t)
	defer teardownDBHelper(t, db)

	t.Run("happy case", func(t *testing.T) {
		// Arrange

		id := uuid.New().String()

		givenUser := &User{
			ID:        id,
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe",
			Password:  "password",
			Email:     "joedoe@foo.bar",
			Country:   "BR",
			CreatedAt: time.Time{}.Add(1 * time.Second),
			UpdatedAt: time.Time{}.Add(2 * time.Second),
		}

		repo := NewPostgres(db)

		// Insert user so we can test the Get method
		require.NoError(t, repo.Insert(context.TODO(), givenUser))

		// Act
		actualUser, actualErr := repo.Get(context.TODO(), id)
		require.NoError(t, actualErr)

		// Assert
		require.Equal(t, givenUser, actualUser)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		repo := NewPostgres(db)

		// Act
		actualUser, actualErr := repo.Get(context.TODO(), uuid.New().String())

		// Assert
		require.Error(t, actualErr)
		assert.True(t, errors.Is(actualErr, ErrUserNotFound))
		assert.Nil(t, actualUser)
	})
}

func TestGetByCountry(t *testing.T) {
	db := setupDBHelper(t)
	defer teardownDBHelper(t, db)

	t.Run("happy case", func(t *testing.T) {
		// Arrange
		givenUsers := []*User{
			{
				ID:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Password:  "password",
				Email:     "joedoe@foo.bar",
				Country:   "BR",
				CreatedAt: time.Time{}.Add(1 * time.Second),
				UpdatedAt: time.Time{}.Add(2 * time.Second),
			},
			{
				ID:        uuid.New().String(),
				FirstName: "Jane",
				LastName:  "Doe",
				Nickname:  "janedoe",
				Password:  "password",
				Email:     "janedoe@foo.bar",
				Country:   "BR",
				CreatedAt: time.Time{}.Add(1 * time.Second),
				UpdatedAt: time.Time{}.Add(2 * time.Second),
			},
		}

		repo := NewPostgres(db)

		// Insert users so we can test the GetByCountry method
		for _, user := range givenUsers {
			require.NoError(t, repo.Insert(context.TODO(), user))
		}

		// Act
		actualUsers, actualErr := repo.GetByCountry(context.TODO(), "BR", "", 10)
		require.NoError(t, actualErr)

		// Assert

		require.Len(t, actualUsers, 2)
		assert.Contains(t, actualUsers, givenUsers[0])
		assert.Contains(t, actualUsers, givenUsers[1])
	})

	t.Run("empty list", func(t *testing.T) {
		// Arrange
		repo := NewPostgres(db)

		// Act
		actualUsers, actualErr := repo.GetByCountry(context.TODO(), "UK", "", 10)

		// Assert
		require.NoError(t, actualErr)
		require.Empty(t, actualUsers)
	})
}

func TestInsert(t *testing.T) {
	db := setupDBHelper(t)
	defer teardownDBHelper(t, db)

	t.Run("happy case", func(t *testing.T) {
		givenUser := &User{
			ID:        uuid.New().String(),
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe",
			Password:  "password",
			Email:     "joedoe@foo.bar",
			Country:   "BR",
			CreatedAt: time.Time{}.Add(1 * time.Second),
			UpdatedAt: time.Time{}.Add(2 * time.Second),
		}

		repo := NewPostgres(db)

		// Act
		actualErr := repo.Insert(context.TODO(), givenUser)
		require.NoError(t, actualErr)

		// Assert
		actualUser, actualErr := repo.Get(context.TODO(), givenUser.ID)
		require.NoError(t, actualErr)

		require.Equal(t, givenUser, actualUser)
	})

	t.Run("duplicate email", func(t *testing.T) {
		// Arrange
		// User was already inserted in the previous test
		repo := NewPostgres(db)

		// Act

		givenUser := &User{
			ID:        uuid.New().String(),
			FirstName: "John",
			LastName:  "Smith",
			Nickname:  "johnsmith",
			Password:  "password",
			Email:     "joedoe@foo.bar",
			Country:   "US",
			CreatedAt: time.Time{}.Add(1 * time.Second),
			UpdatedAt: time.Time{}.Add(2 * time.Second),
		}

		actualErr := repo.Insert(context.TODO(), givenUser)

		// Assert
		require.Error(t, actualErr)
		assert.True(t, errors.Is(actualErr, ErrDuplicateEmail))
	})
}

func TestUpdate(t *testing.T) {
	db := setupDBHelper(t)
	defer teardownDBHelper(t, db)

	t.Run("happy case", func(t *testing.T) {
		// Arrange
		givenUser := &User{
			ID:        uuid.New().String(),
			FirstName: "John",
			LastName:  "Doe",
			Nickname:  "johndoe",
			Password:  "password",
			Email:     "joedoe@foo.bar",
			Country:   "BR",
			CreatedAt: time.Time{}.Add(1 * time.Second),
			UpdatedAt: time.Time{}.Add(2 * time.Second),
		}

		repo := NewPostgres(db)

		// Insert a user manually so we can test the Update method
		require.NoError(t, repo.Insert(context.TODO(), givenUser))

		// Act
		err := repo.Update(context.TODO(), &User{
			ID:        givenUser.ID,
			FirstName: "Joe",
			LastName:  "Doe",
			Nickname:  "hollywoodjoe",
			Password:  "password",
			Email:     "joedoe@foo.quz",
			Country:   "US",
			CreatedAt: time.Time{}.Add(1 * time.Second),
			UpdatedAt: time.Time{}.Add(2 * time.Second),
		})
		require.NoError(t, err)

		// Assert

		actualUser, actualErr := repo.Get(context.TODO(), givenUser.ID)
		require.NoError(t, actualErr)

		require.Equal(t, "Joe", actualUser.FirstName)
		require.Equal(t, "Doe", actualUser.LastName)
		require.Equal(t, "hollywoodjoe", actualUser.Nickname)
		require.Equal(t, "password", actualUser.Password)
		require.Equal(t, "joedoe@foo.quz", actualUser.Email)
		require.Equal(t, "US", actualUser.Country)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		repo := NewPostgres(db)

		// Act
		err := repo.Update(context.TODO(), &User{
			ID:        uuid.New().String(),
			FirstName: "Joe",
			LastName:  "Doe",
			Nickname:  "hollywoodjoe",
			Password:  "password",
			Email:     "joedoe@foo.quz",
			Country:   "US",
			CreatedAt: time.Time{}.Add(1 * time.Second),
			UpdatedAt: time.Time{}.Add(2 * time.Second),
		})

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUserNotFound))
	})
}

const (
	migrationsDir      string = "../../../migrations"
	postgresDriverName string = "postgres"
	dbHost             string = "localhost"
	dbPort             string = "5432"
	dbUser             string = "user"
	dbPass             string = "password"
	dbName             string = "usrsvc"
)

func setupDBHelper(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open(postgresDriverName, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName),
	)
	require.NoError(t, err)

	require.NoError(t, goose.Up(db.DB, migrationsDir))
	return db
}

func teardownDBHelper(t *testing.T, db *sqlx.DB) {
	t.Helper()

	_, err := db.Exec("TRUNCATE TABLE users CASCADE")
	require.NoError(t, err)

	require.NoError(t, db.Close())
}
