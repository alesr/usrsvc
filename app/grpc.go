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
	defaultCursor   string        = ""
)

type userService interface {
	Fetch(ctx context.Context, id string) (*service.User, error)
	FetchAll(ctx context.Context, filter service.FilterParams, pag service.PaginationParams) ([]*service.User, error)
	Create(ctx context.Context, user *service.User) (*service.User, error)
	Update(ctx context.Context, user *service.User) (*service.User, error)
	Delete(ctx context.Context, id string) error
}

type GRPCServer struct {
	apiv1.UnimplementedUserServiceServer
	logger  *zap.Logger
	service userService
}

func NewGRPCServer(logger *zap.Logger, service userService) *GRPCServer {
	return &GRPCServer{
		logger:  logger,
		service: service,
	}
}

// Register the GRPCServer to a gRPC server
func (s *GRPCServer) Register(server *grpc.Server) {
	apiv1.RegisterUserServiceServer(server, s)
	log.Println("Registered GRPCServer to gRPC server")
}

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

func newUserResponseFromDomain(user *service.User) *apiv1.User {
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
