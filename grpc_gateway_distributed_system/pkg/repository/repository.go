package repository

import (
	"database/sql"

	author_client "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/client/author"
	author_repository "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository/author"
	book_repository "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository/book"
	author_repository_impl "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repositoryimpl/author"
	book_repository_impl "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repositoryimpl/book"
)

type RepositorySet struct {
	Author       author_repository.AuthorRepository
	AuthorClient author_client.Client
	Book         book_repository.BookRepository
}

func SetUp(db *sql.DB) RepositorySet {
	// gRPCサーバーに接続
	authorGRPCClient, err := author_client.NewGRPCClient("author-service:9092")
	if err != nil {
		// gRPC接続に失敗した場合はパニック（HTTPフォールバックは削除）
		panic("Failed to connect to author gRPC service: " + err.Error())
	}

	return RepositorySet{
		Author:       author_repository_impl.NewAuthorRepository(db),
		AuthorClient: authorGRPCClient,
		Book:         book_repository_impl.NewBookRepository(db),
	}
}
