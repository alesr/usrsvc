package app

import (
	"errors"
	"testing"

	apiv1 "github.com/alesr/usrsvc/proto/users/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateUserRequest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		given    *apiv1.CreateUserRequest
		expected error
	}{
		{
			name: "valid request",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: nil,
		},
		{
			name: "invalid first name",
			given: &apiv1.CreateUserRequest{
				FirstName: "J",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameLength,
		},
		{
			name: "missing first name",
			given: &apiv1.CreateUserRequest{
				FirstName: "",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameRequired,
		},
		{
			name: "invalid last name",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "D",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameLength,
		},
		{
			name: "missing last name",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameRequired,
		},
		{
			name: "invalid nickname",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "j",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameLength,
		},
		{
			name: "missing nickname",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameRequired,
		},
		{
			name: "invalid email",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrEmailFormat,
		},
		{
			name: "missing email",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrEmailRequired,
		},
		{
			name: "invalid password",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "invalidpassword",
				Country:   "BR",
			},
			expected: ErrPasswordFormat,
		},
		{
			name: "missing password",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "",
				Country:   "BR",
			},
			expected: ErrPasswordRequired,
		},
		{
			name: "password length",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "xxx",
				Country:   "BR",
			},
			expected: ErrPasswordLength,
		},
		{
			name: "invalid country code",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BRR",
			},
			expected: ErrCountryCodeInvalid,
		},
		{
			name: "missing country code",
			given: &apiv1.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "",
			},
			expected: ErrCountryCodeRequired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			observedErr := validateCreateUserRequest(tc.given)
			assert.True(t, errors.Is(observedErr, tc.expected))
		})
	}
}

func TestValidateUpdateUserRequest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		given    *apiv1.UpdateUserRequest
		expected error
	}{
		{
			name: "valid request",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: nil,
		},
		{
			name: "invalid id",
			given: &apiv1.UpdateUserRequest{
				Id:        "123",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrIDFormat,
		},
		{
			name: "missing id",
			given: &apiv1.UpdateUserRequest{
				Id:        "",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrIDRequired,
		},
		{
			name: "invalid first name",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "J",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameLength,
		},
		{
			name: "missing first name",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameRequired,
		},
		{
			name: "invalid last name",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "D",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameLength,
		},
		{
			name: "missing last name",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameRequired,
		},
		{
			name: "invalid nickname",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "j",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameLength,
		},
		{
			name: "missing nickname",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrNameRequired,
		},
		{
			name: "invalid email",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "foo.bar",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrEmailFormat,
		},
		{
			name: "missing email",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "",
				Password:  "some_passw0rd",
				Country:   "BR",
			},
			expected: ErrEmailRequired,
		},
		{
			name: "invalid password",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "invalidpassword",
				Country:   "BR",
			},
			expected: ErrPasswordFormat,
		},
		{
			name: "missing password",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "",
				Country:   "BR",
			},
			expected: ErrPasswordRequired,
		},
		{
			name: "password length",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "xxx",
				Country:   "BR",
			},
			expected: ErrPasswordLength,
		},
		{
			name: "invalid country code",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "BRR",
			},
			expected: ErrCountryCodeInvalid,
		},
		{
			name: "missing country code",
			given: &apiv1.UpdateUserRequest{
				Id:        uuid.New().String(),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "johndoe",
				Email:     "joedoe@foo.bar",
				Password:  "some_passw0rd",
				Country:   "",
			},
			expected: ErrCountryCodeRequired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			observedErr := validateUpdateUserRequest(tc.given)
			assert.True(t, errors.Is(observedErr, tc.expected))
		})
	}
}
