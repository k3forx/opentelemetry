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
	author_client "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/client/author"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/grpc/server"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository"
	author_repository_impl "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repositoryimpl/author"
	book_repository_impl "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repositoryimpl/book"
	bookv1 "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/proto/gen/book/v1"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
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

	repositorySet := setupBookRepository(db)

	// BookのgRPCサーバーを起動
	if err := startBookGRPCServer(repositorySet); err != nil {
		log.Fatal("Failed to start book gRPC server:", err)
	}
}

func setupBookRepository(db *sql.DB) repository.RepositorySet {
	// Author ServiceのgRPCクライアント
	authorGRPCClient, err := author_client.NewGRPCClient("author-service:9092")
	if err != nil {
		log.Fatal("Failed to connect to author service:", err)
	}

	return repository.RepositorySet{
		Author:       author_repository_impl.NewAuthorRepository(db),
		AuthorClient: authorGRPCClient,
		Book:         book_repository_impl.NewBookRepository(db),
	}
}

func startBookGRPCServer(repositorySet repository.RepositorySet) error {
	port := os.Getenv("BOOK_GRPC_PORT")
	if port == "" {
		port = "9091"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	bookServer := server.NewBookServer(repositorySet)
	bookv1.RegisterBookServiceServer(s, bookServer)

	log.Printf("Book gRPC server listening on :%s", port)
	return s.Serve(lis)
}
