package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/otel/trace"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/grpc/server"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	authorv1 "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/proto/gen/author/v1"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository"
	author_repository_impl "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repositoryimpl/author"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	tp, err := trace.InitTraceProvider(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root_password",
		DBName:               "app",
		Addr:                 "mysql-grpc-gateway:3306",
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	
	// Open database with OTel SQL instrumentation
	db, err := otelsql.Open("mysql", cfg.FormatDSN(),
		otelsql.WithAttributes(semconv.DBSystemMySQL),
		otelsql.WithDBName(cfg.DBName),
	)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	repositorySet := setupAuthorRepository(db)

	// AuthorのgRPCサーバーを起動
	if err := startAuthorGRPCServer(repositorySet); err != nil {
		log.Fatal("Failed to start author gRPC server:", err)
	}
}

func setupAuthorRepository(db *sql.DB) repository.RepositorySet {
	return repository.RepositorySet{
		Author: author_repository_impl.NewAuthorRepository(db),
	}
}

func startAuthorGRPCServer(repositorySet repository.RepositorySet) error {
	port := os.Getenv("AUTHOR_GRPC_PORT")
	if port == "" {
		port = "9092"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	authorServer := server.NewAuthorServer(repositorySet)
	authorv1.RegisterAuthorServiceServer(s, authorServer)

	log.Printf("Author gRPC server listening on :%s", port)
	return s.Serve(lis)
}
