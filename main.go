package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/alesr/usrsvc/app"
	userrepo "github.com/alesr/usrsvc/internal/users/repository"
	userservice "github.com/alesr/usrsvc/internal/users/service"
	"github.com/alesr/usrsvc/pkg/events"
	apiv1 "github.com/alesr/usrsvc/proto/users/v1"
	"github.com/jmoiron/sqlx"
	envars "github.com/netflix/go-env"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const (
	postgresDriverName string = "postgres"
	dbMigrationsDir    string = "migrations"
	grpcPort           string = ":50051"
)

type config struct {
	DBUser string `env:"POSTGRES_USER,default=user"`
	DBPass string `env:"POSTGRES_PASSWORD,default=password"`
	DBName string `env:"POSTGRES_DB,default=usrsvc"`
	DBHost string `env:"POSTGRES_HOST,default=db"`
	DBPort string `env:"POSTGRES_PORT,default=5432"`
}

func newConfig() *config {
	var cfg config
	if _, err := envars.UnmarshalFromEnviron(&cfg); err != nil {
		log.Fatal(err)
	}
	return &cfg
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("failed to create logger", err)
	}

	defer logger.Sync()

	cfg := newConfig()

	db, err := sqlx.Open(postgresDriverName, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName),
	)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(postgresDriverName); err != nil {
		logger.Fatal("failed to set goose dialect", zap.Error(err))
	}

	if err := goose.Up(db.DB, dbMigrationsDir); err != nil {
		logger.Fatal("failed to run goose migrations", zap.Error(err))
	}

	userRepo := userrepo.NewPostgres(db)

	userService := userservice.NewServiceDefault(
		logger,
		userRepo,
		userservice.WithPublisher(&fakePubSub{}),
	)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.Fatal("failed to listen on grpc port", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	grpcServer.RegisterService(
		&apiv1.UserService_ServiceDesc,
		app.NewGRPCServer(logger, userService),
	)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve gRPC server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer signal.Stop(c)

	<-c
	logger.Info("shutting down gRPC server")
	grpcServer.GracefulStop()
}

type fakePubSub struct{}

func (f *fakePubSub) Publish(event events.Event, data any) error {
	return nil
}
