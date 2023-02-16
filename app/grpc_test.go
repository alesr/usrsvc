package app

import (
	"testing"
	"time"

	"github.com/alesr/usrsvc/internal/users/service"
	apiv1 "github.com/alesr/usrsvc/proto/users/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
