package app

import (
	"context"
	"testing"
	"time"

	"github.com/alesr/usrsvc/internal/users/service"
	apiv1 "github.com/alesr/usrsvc/proto/users/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()
	createdAt := time.Time{}.Add(1 * time.Second)
	updatedAt := time.Time{}.Add(2 * time.Second)

	var createFuncWasCalled bool
	svc := &serviceMock{
		CreateFunc: func(ctx context.Context, user *service.User) (*service.User, error) {
			createFuncWasCalled = true
			return &service.User{
				ID:        id,
				FirstName: "Michael",
				LastName:  "Jackson",
				Nickname:  "mj",
				Email:     "mj@foo.bar",
				Country:   "US",
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			}, nil
		},
	}

	server := NewGRPCServer(zap.NewNop(), svc)

	givenReq := &apiv1.CreateUserRequest{
		FirstName: "Michael",
		LastName:  "Jackson",
		Nickname:  "mj",
		Email:     "mj@foo.bar",
		Password:  "some-passw0rd",
		Country:   "US",
	}

	expectedResp := &apiv1.CreateUserResponse{
		User: &apiv1.User{
			Id:        id,
			FirstName: "Michael",
			LastName:  "Jackson",
			Nickname:  "mj",
			Email:     "mj@foo.bar",
			Country:   "US",
		},
	}

	observed, err := server.CreateUser(context.TODO(), givenReq)
	assert.NoError(t, err)

	assert.True(t, createFuncWasCalled)

	assert.Equal(t, expectedResp.User.Id, observed.User.Id)
	assert.Equal(t, expectedResp.User.FirstName, observed.User.FirstName)
	assert.Equal(t, expectedResp.User.LastName, observed.User.LastName)
	assert.Equal(t, expectedResp.User.Nickname, observed.User.Nickname)
	assert.Equal(t, expectedResp.User.Email, observed.User.Email)
	assert.Equal(t, expectedResp.User.Country, observed.User.Country)
	assert.Equal(t, timestamppb.New(createdAt), observed.User.CreatedAt)
	assert.Equal(t, timestamppb.New(updatedAt), observed.User.UpdatedAt)

	t.Run("when the service returns an error", func(t *testing.T) {
		t.SkipNow()
		// We alread have good test coverage, and we got the idea...
		// I'll leave these tests as an exercise for the reader =]
	})

	t.Run("when the request is invalid", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the email is already in use", func(t *testing.T) {
		t.SkipNow()
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the service returns an error", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the user is not found", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the request is invalid", func(t *testing.T) {
		t.SkipNow()
	})
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the service returns an error", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the user is not found", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the request is invalid", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the email is already in use", func(t *testing.T) {
		t.SkipNow()
	})
}

func TestListUser(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the service returns an error", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the request is invalid", func(t *testing.T) {
		t.SkipNow()
	})
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the service returns an error", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the user is not found", func(t *testing.T) {
		t.SkipNow()
	})

	t.Run("when the request is invalid", func(t *testing.T) {
		t.SkipNow()
	})
}

func TestNewUserResponseFromDomain(t *testing.T) {
	t.Parallel()

	id := uuid.New().String()

	testCases := []struct {
		name     string
		given    *service.User
		expected *apiv1.User
	}{
		{
			name: "happy path",
			given: &service.User{
				ID:        id,
				FirstName: "Michael",
				LastName:  "Jackson",
				Nickname:  "mj",
				Email:     "mj@foo.bar",
				Country:   "US",
				CreatedAt: time.Time{}.Add(1 * time.Second),
				UpdatedAt: time.Time{}.Add(2 * time.Second),
			},
			expected: &apiv1.User{
				Id:        id,
				FirstName: "Michael",
				LastName:  "Jackson",
				Nickname:  "mj",
				Email:     "mj@foo.bar",
				Country:   "US",
				CreatedAt: timestamppb.New(time.Time{}.Add(1 * time.Second)),
				UpdatedAt: timestamppb.New(time.Time{}.Add(2 * time.Second)),
			},
		},
		{
			name:     "nil user",
			given:    nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			observed := newUserResponseFromDomain(tc.given)
			assert.Equal(t, tc.expected, observed)
		})
	}
}
