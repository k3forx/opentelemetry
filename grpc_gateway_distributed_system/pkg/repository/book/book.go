package book_repository

import (
	"context"

	book_model "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/model/book"
)

type BookRepository interface {
	GetByID(ctx context.Context, id int64) (book_model.Book, error)
}
