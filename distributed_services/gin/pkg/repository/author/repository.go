package author_repository

import (
	"context"

	author_model "github.com/k3forx/opentelemetry/gin/pkg/model/author"
)

type AuthorRepository interface {
	Create(ctx context.Context, author *author_model.Author) error
	GetByID(ctx context.Context, id int64) (author_model.Author, error)
}
