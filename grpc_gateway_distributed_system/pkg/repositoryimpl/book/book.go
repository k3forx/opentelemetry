package book_repository_impl

import (
	"context"
	"database/sql"
	"errors"

	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/otel/trace"
	book_model "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/model/book"
	book_repository "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository/book"
)

var _ book_repository.BookRepository = &bookRepositoryImpl{}

type bookRepositoryImpl struct {
	queries *Queries
}

func NewBookRepository(db *sql.DB) *bookRepositoryImpl {
	return &bookRepositoryImpl{queries: New(db)}
}

func (impl *bookRepositoryImpl) GetByID(ctx context.Context, id int64) (book_model.Book, error) {
	ctx, span := trace.Tracer.Start(ctx, trace.SpanNameRepository)
	defer span.End()

	b, err := impl.queries.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return book_model.Book{}, nil
		}
	}
	return book_model.Book{ID: b.ID, Title: b.Title, AuthorID: b.AuthorID}, nil
}
