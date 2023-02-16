package tests

import (
	"fmt"
	"net"
	"testing"

	"github.com/alesr/usrsvc/app"
	"github.com/alesr/usrsvc/internal/users/repository"
	"github.com/alesr/usrsvc/internal/users/service"
	apiv1 "github.com/alesr/usrsvc/proto/users/v1"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	grpcPort           string = ":50051"
	migrationsDir      string = "../migrations"
	postgresDriverName string = "postgres"
	dbHost             string = "localhost"
	dbPort             string = "5432"
	dbUser             string = "user"
	dbPass             string = "password"
	dbName             string = "usrsvc"
)

func startGRPCServerHelper(t *testing.T, db *sqlx.DB) func() {
	t.Helper()

	grpcServer := grpc.NewServer()

	grpcServer.RegisterService(
		&apiv1.UserService_ServiceDesc,
		app.NewGRPCServer(
			zap.NewNop(),
			service.NewServiceDefault(
				zap.NewNop(),
				repository.NewPostgres(db),
			),
		),
	)

	lis, err := net.Listen("tcp", grpcPort)
	require.NoError(t, err)

	go func() {
		err := grpcServer.Serve(lis)
		require.NoError(t, err)
	}()

	return grpcServer.Stop
}

func setupGRPClientHelper(t *testing.T) (apiv1.UserServiceClient, func() error) {
	t.Helper()

	conn, err := grpc.Dial(grpcPort, grpc.WithInsecure())
	require.NoError(t, err)

	return apiv1.NewUserServiceClient(conn), conn.Close
}

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
