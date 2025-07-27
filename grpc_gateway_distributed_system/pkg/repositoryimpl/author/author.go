package author_repository_impl

import (
	"context"
	"database/sql"
	"errors"

	author_model "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/model/author"
	author_repository "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository/author"
)

var _ author_repository.AuthorRepository = &authorRepositoryImpl{}

type authorRepositoryImpl struct {
	queries *Queries
}

func NewAuthorRepository(db *sql.DB) *authorRepositoryImpl {
	return &authorRepositoryImpl{queries: New(db)}
}

func (impl *authorRepositoryImpl) GetByID(ctx context.Context, id int64) (author_model.Author, error) {
	// ctx, span := trace.Tracer.Start(ctx, trace.SpanNameRepository)
	// defer span.End()

	a, err := impl.queries.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return author_model.Author{}, nil
		}
		return author_model.Author{}, err
	}
	return author_model.Author{ID: id, Name: a.Name, Bio: a.Bio}, nil
}
