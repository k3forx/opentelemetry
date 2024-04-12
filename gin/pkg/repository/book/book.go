package book_repository

import (
	"context"

	book_model "github.com/k3forx/opentelemetry/gin/pkg/model/book"
)

type BookRepository interface {
	GetAllByAuthorID(ctx context.Context, id int64) ([]book_model.Book, error)
	GetByID(ctx context.Context, id int64) (book_model.Book, error)
}
