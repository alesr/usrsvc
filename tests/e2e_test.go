//go:build e2e
// +build e2e

package tests

import (
	"context"
	"errors"
	"testing"

	apiv1 "github.com/alesr/usrsvc/proto/users/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_E2E(t *testing.T) {
	db := setupDBHelper(t)
	defer teardownDBHelper(t, db)

	stopServer := startGRPCServerHelper(t, db)
	defer stopServer()

	grpcClient, close := setupGRPClientHelper(t)
	defer func() {
		err := close()
		require.NoError(t, err)
	}()

	// First, create a user

	givenCreateReq := &apiv1.CreateUserRequest{
		FirstName: "Michael",
		LastName:  "Jackson",
		Nickname:  "mj",
		Email:     "mj@foo.bar",
		Password:  "s0meP@ssw0rd",
		Country:   "US",
	}

	// I took the decision of not returning the password in (all) responses.
	// This is a security measure, as the password should never be sent to the client.
	// It's also stored as a hash, so it's not possible to retrieve it =]
	expectedCreateResp := &apiv1.CreateUserResponse{
		User: &apiv1.User{
			FirstName: "Michael",
			LastName:  "Jackson",
			Nickname:  "mj",
			Email:     "mj@foo.bar",
			Country:   "US",
		},
	}

	observedCreateResp, err := grpcClient.CreateUser(context.TODO(), givenCreateReq)
	require.NoError(t, err)

	assert.NotEmpty(t, observedCreateResp.User.Id)
	assert.NotEmpty(t, observedCreateResp.User.CreatedAt)
	assert.NotEmpty(t, observedCreateResp.User.UpdatedAt)

	assert.Equal(t, expectedCreateResp.User.FirstName, observedCreateResp.User.FirstName)
	assert.Equal(t, expectedCreateResp.User.LastName, observedCreateResp.User.LastName)
	assert.Equal(t, expectedCreateResp.User.Nickname, observedCreateResp.User.Nickname)
	assert.Equal(t, expectedCreateResp.User.Email, observedCreateResp.User.Email)
	assert.Equal(t, expectedCreateResp.User.Country, observedCreateResp.User.Country)

	// Second, get the user

	givenGetReq := &apiv1.GetUserRequest{Id: observedCreateResp.User.Id}

	observedGetResp, err := grpcClient.GetUser(context.TODO(), givenGetReq)
	require.NoError(t, err)

	assert.Equal(t, observedCreateResp.User.Id, observedGetResp.User.Id)
	assert.Equal(t, observedCreateResp.User.FirstName, observedGetResp.User.FirstName)
	assert.Equal(t, observedCreateResp.User.LastName, observedGetResp.User.LastName)
	assert.Equal(t, observedCreateResp.User.Nickname, observedGetResp.User.Nickname)
	assert.Equal(t, observedCreateResp.User.Email, observedGetResp.User.Email)
	assert.Equal(t, observedCreateResp.User.Country, observedGetResp.User.Country)
	assert.NotEmpty(t, observedGetResp.User.CreatedAt)
	assert.NotEmpty(t, observedGetResp.User.UpdatedAt)

	// Third, we get our user in the list

	givenListReq := &apiv1.ListUsersRequest{}

	observedListResp, err := grpcClient.ListUsers(context.TODO(), givenListReq)
	require.NoError(t, err)

	require.Len(t, observedListResp.Users, 1)

	assert.Equal(t, observedCreateResp.User.Id, observedListResp.Users[0].Id)
	assert.Equal(t, observedCreateResp.User.FirstName, observedListResp.Users[0].FirstName)
	assert.Equal(t, observedCreateResp.User.LastName, observedListResp.Users[0].LastName)
	assert.Equal(t, observedCreateResp.User.Nickname, observedListResp.Users[0].Nickname)
	assert.Equal(t, observedCreateResp.User.Email, observedListResp.Users[0].Email)
	assert.Equal(t, observedCreateResp.User.Country, observedListResp.Users[0].Country)
	assert.NotEmpty(t, observedListResp.Users[0].CreatedAt)
	assert.NotEmpty(t, observedListResp.Users[0].UpdatedAt)

	// Fourth, we also get our user in the list when filtering by country

	givenListReq = &apiv1.ListUsersRequest{Country: "US"}

	observedListFilteredResp, err := grpcClient.ListUsers(context.TODO(), givenListReq)
	require.NoError(t, err)

	require.Len(t, observedListFilteredResp.Users, 1)

	assert.Equal(t, observedCreateResp.User.Id, observedListFilteredResp.Users[0].Id)
	assert.Equal(t, observedCreateResp.User.FirstName, observedListFilteredResp.Users[0].FirstName)
	assert.Equal(t, observedCreateResp.User.LastName, observedListFilteredResp.Users[0].LastName)
	assert.Equal(t, observedCreateResp.User.Nickname, observedListFilteredResp.Users[0].Nickname)
	assert.Equal(t, observedCreateResp.User.Email, observedListFilteredResp.Users[0].Email)
	assert.Equal(t, observedCreateResp.User.Country, observedListFilteredResp.Users[0].Country)
	assert.NotEmpty(t, observedListFilteredResp.Users[0].CreatedAt)
	assert.NotEmpty(t, observedListFilteredResp.Users[0].UpdatedAt)

	// Fifth, we update the user and check that the changes are reflected

	// In a real world scenario, I probably wouldn't allow
	// the user to update the email, but for the sake of keeping it simple...

	givenUpdateReq := &apiv1.UpdateUserRequest{
		Id:        observedCreateResp.User.Id,
		FirstName: "Magic",
		LastName:  "Jordan",
		Nickname:  "magic",
		Email:     "magic@foo.bar",
		Password:  "s0meP@ssw0rd2",
		Country:   "BR",
	}

	expectedUpdateResp := &apiv1.UpdateUserResponse{
		User: &apiv1.User{
			Id:        observedCreateResp.User.Id,
			FirstName: "Magic",
			LastName:  "Jordan",
			Nickname:  "magic",
			Email:     "magic@foo.bar",
			Country:   "BR",
		},
	}

	observedUpdateResp, err := grpcClient.UpdateUser(context.TODO(), givenUpdateReq)
	require.NoError(t, err)

	assert.Equal(t, expectedUpdateResp.User.FirstName, observedUpdateResp.User.FirstName)
	assert.Equal(t, expectedUpdateResp.User.LastName, observedUpdateResp.User.LastName)
	assert.Equal(t, expectedUpdateResp.User.Nickname, observedUpdateResp.User.Nickname)
	assert.Equal(t, expectedUpdateResp.User.Email, observedUpdateResp.User.Email)
	assert.Equal(t, expectedUpdateResp.User.Country, observedUpdateResp.User.Country)
	assert.NotEmpty(t, observedUpdateResp.User.CreatedAt)
	assert.NotEmpty(t, observedUpdateResp.User.UpdatedAt)

	// Finally, we delete the user and check that it's no longer returned =[

	givenDeleteReq := &apiv1.DeleteUserRequest{Id: observedCreateResp.User.Id}

	_, err = grpcClient.DeleteUser(context.TODO(), givenDeleteReq)
	require.NoError(t, err)

	givenGetReq = &apiv1.GetUserRequest{Id: observedCreateResp.User.Id}

	observedGetResp, err = grpcClient.GetUser(context.TODO(), givenGetReq)
	require.Error(t, err)

	assert.True(t, errors.Is(err, status.Error(codes.NotFound, "user not found")))
	assert.Nil(t, observedGetResp)
}
