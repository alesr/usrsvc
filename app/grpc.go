package app

import (
	"context"
	"log"
	"time"

	"github.com/alesr/usrsvc/internal/users/service"
	apiv1 "github.com/alesr/usrsvc/proto/users/v1"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	ctxTimeout      time.Duration = 5 * time.Second
	defaultPageSize int32         = 100
)

// userService is the interface that provides the business logic for the gRPC server.
type userService interface {
	Fetch(ctx context.Context, id string) (*service.User, error)
	FetchAll(ctx context.Context, filter service.FilterParams, pag service.PaginationParams) ([]*service.User, error)
	Create(ctx context.Context, user *service.User) (*service.User, error)
	Update(ctx context.Context, user *service.User) (*service.User, error)
	Delete(ctx context.Context, id string) error
	CheckServiceHealth(ctx context.Context) error
}

// GRPCServer is the gRPC server that provides the user service.
type GRPCServer struct {
	apiv1.UnimplementedUserServiceServer
	logger  *zap.Logger
	service userService
}

// NewGRPCServer creates a new gRPC server.
func NewGRPCServer(logger *zap.Logger, service userService) *GRPCServer {
	return &GRPCServer{
		logger:  logger,
		service: service,
	}
}

// Register registers the gRPC server to (our) GRPCServer.
func (s *GRPCServer) Register(server *grpc.Server) {
	apiv1.RegisterUserServiceServer(server, s)
	log.Println("Registered GRPCServer to gRPC server")
}

// GetUser returns a user by ID.
func (s *GRPCServer) CreateUser(ctx context.Context, req *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
	if err := validateCreateUserRequest(req); err != nil {
		s.logger.Error("failed to validate request", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	user := &service.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Password:  req.Password,
		Country:   req.Country,
	}

	user, err := s.service.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, convertServiceError(err)
	}

	return &apiv1.CreateUserResponse{
		User: newUserResponseFromDomain(user),
	}, nil
}

// UpdateUser updates a user by ID.
// For the sake of simplicity, we update all the fields of the user but the ID.
func (s *GRPCServer) UpdateUser(ctx context.Context, req *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error) {
	if err := validateUpdateUserRequest(req); err != nil {
		s.logger.Error("failed to validate request", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	user := &service.User{
		ID:        req.Id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Password:  req.Password,
		Country:   req.Country,
	}

	user, err := s.service.Update(ctx, user)
	if err != nil {
		s.logger.Error("failed to update user", zap.Error(err))
		return nil, convertServiceError(err)
	}

	return &apiv1.UpdateUserResponse{
		User: newUserResponseFromDomain(user),
	}, nil
}

// DeleteUser deletes a user by ID.
func (s *GRPCServer) GetUser(ctx context.Context, req *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	if err := validateID(req.Id); err != nil {
		s.logger.Error("failed to validate id", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	user, err := s.service.Fetch(ctx, req.Id)
	if err != nil {
		s.logger.Error("failed to fetch user", zap.Error(err))
		return nil, convertServiceError(err)
	}

	return &apiv1.GetUserResponse{
		User: newUserResponseFromDomain(user),
	}, nil
}

/*
ListUsers returns a list of users.

The list can be filtered by country and paginated.
The default page size is 100. If a page size is not provided or is invalid, the default page size is used.
The default page token points to the last ID in the list.
If a page token is not required, but if an invalid page token is provided, an error is returned.

The implementation for the pagination is based on https://cloud.google.com/apis/design/design_patterns#list_pagination
*/
func (s *GRPCServer) ListUsers(ctx context.Context, req *apiv1.ListUsersRequest) (*apiv1.ListUsersResponse, error) {
	if req.PageSize <= 0 || req.PageSize > defaultPageSize {
		req.PageSize = defaultPageSize
	}

	if req.PageToken != "" {
		if _, err := uuid.Parse(req.PageToken); err != nil {
			s.logger.Error("failed to validate cursor", zap.Error(err))
			return nil, ErrPageTokenInvalid
		}
	}

	if req.Country != "" && len(req.Country) != 2 {
		return nil, ErrCountryCodeInvalid
	}

	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	var filters service.FilterParams
	if req.Country != "" {
		filters.Country = &req.Country
	}

	pagination := service.PaginationParams{
		Limit:  int(req.PageSize),
		Cursor: req.PageToken,
	}

	users, err := s.service.FetchAll(ctx, filters, pagination)
	if err != nil {
		s.logger.Error("failed to fetch users", zap.Error(err))
		return nil, convertServiceError(err)
	}

	var nextPageToken string
	if len(users) == int(req.PageSize) {
		nextPageToken = users[len(users)-1].ID
	}

	var usersProto []*apiv1.User
	for _, user := range users {
		usersProto = append(usersProto, newUserResponseFromDomain(user))
	}

	return &apiv1.ListUsersResponse{
		Users:         usersProto,
		NextPageToken: nextPageToken,
	}, nil
}

// DeleteUser deletes a user by ID.
func (s *GRPCServer) DeleteUser(ctx context.Context, req *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	if err := validateID(req.Id); err != nil {
		s.logger.Error("failed to validate id", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := s.service.Delete(ctx, req.Id); err != nil {
		s.logger.Error("failed to delete user", zap.Error(err))
		return nil, convertServiceError(err)
	}

	return &apiv1.DeleteUserResponse{}, nil
}

// CheckHeath checks the health of the application going all the way down to the database.
func (s *GRPCServer) CheckHeath(ctx context.Context, req *apiv1.HealthCheckRequest) (*apiv1.HealthCheckResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	if err := s.service.CheckServiceHealth(ctx); err != nil {
		return &apiv1.HealthCheckResponse{
			Status: apiv1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	return &apiv1.HealthCheckResponse{
		Status: apiv1.HealthCheckResponse_SERVING,
	}, nil
}

func newUserResponseFromDomain(user *service.User) *apiv1.User {
	// Better safe than sorry.
	if user == nil {
		return nil
	}

	return &apiv1.User{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
