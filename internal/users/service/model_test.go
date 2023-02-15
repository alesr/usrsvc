package service

import (
	"testing"
	"time"

	"github.com/alesr/usrsvc/internal/users/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewUserFromDomain(t *testing.T) {
	given := &User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  "password",
		Email:     "joedoe@foo.bar",
		Country:   "BR",
		CreatedAt: time.Time{}.Add(1 * time.Hour),
		UpdatedAt: time.Time{}.Add(2 * time.Hour),
	}

	expected := &repository.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  "password",
		Email:     "joedoe@foo.bar",
		Country:   "BR",
		CreatedAt: time.Time{}.Add(1 * time.Hour),
		UpdatedAt: time.Time{}.Add(2 * time.Hour),
	}

	actual := newUserStoreFromDomain(given)
	assert.Equal(t, expected, actual)
}

func TestNewUserFromStore(t *testing.T) {
	given := &repository.User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  "password",
		Email:     "joedoe@foo.bar",
		Country:   "BR",
		CreatedAt: time.Time{}.Add(1 * time.Hour),
		UpdatedAt: time.Time{}.Add(2 * time.Hour),
	}

	expected := &User{
		ID:        "123",
		FirstName: "John",
		LastName:  "Doe",
		Nickname:  "johndoe",
		Password:  "password",
		Email:     "joedoe@foo.bar",
		Country:   "BR",
		CreatedAt: time.Time{}.Add(1 * time.Hour),
		UpdatedAt: time.Time{}.Add(2 * time.Hour),
	}

	actual := newUserDomainFromStore(given)
	assert.Equal(t, expected, actual)
}
