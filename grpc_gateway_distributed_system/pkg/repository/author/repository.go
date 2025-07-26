package author_repository

import (
	"context"

	author_model "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/model/author"
)

type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (author_model.Author, error)
}
