package main

import (
	"context"
	"database/sql"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "usdt/internal/handler/pb"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"usdt/internal/client"
	"usdt/internal/config"
	"usdt/internal/handler"
	"usdt/internal/repository"
	"usdt/internal/service"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
	defer logger.Sync()

	tp, err := initTracer()
	if err != nil {
		logger.Fatal("failed to initialize tracer", zap.Error(err))
	}
	defer func() { _ = tp.Shutdown(context.Background()) }()

	cfg := config.Load()
	logger.Info("configuration loaded", zap.String("db_url", cfg.DBUrl), zap.String("port", cfg.Port))

	if err := runMigrations(cfg.DBUrl, logger); err != nil {
		logger.Fatal("migration failed", zap.Error(err))
	}
	logger.Info("database migrations completed")

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		logger.Fatal("failed to open database connection", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("error closing DB", zap.Error(err))
		}
	}()
	logger.Info("database connection established")

	// Initialize repository with logger and tracer
	repo := &repository.Repo{DB: db, Tracer: otel.Tracer("usdt/repository"), Logger: logger}
	client := &client.GrinexClient{URL: cfg.GrinexAPIUrl}
	svc := &service.RateService{Client: client, Repo: repo}

	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logger.Fatal("failed to listen on port", zap.String("port", cfg.Port), zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRateServiceServer(grpcServer, &handler.Server{Service: svc})

	go func() {
		logger.Info("gRPC server started", zap.String("port", cfg.Port))
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	sigReceived := <-sig

	logger.Info("signal received, shutting down gracefully", zap.String("signal", sigReceived.String()))
	grpcServer.GracefulStop()
}

func runMigrations(dbURL string, logger *zap.Logger) error {
	m, err := migrate.New("file:///migrations", dbURL)
	if err != nil {
		logger.Error("failed to create migration instance", zap.Error(err))
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error("migration failed", zap.Error(err))
		return err
	}
	return nil
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("usdt-app"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
