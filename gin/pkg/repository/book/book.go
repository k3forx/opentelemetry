package book_repository

import (
	"context"

	book_model "github.com/k3forx/opentelemetry/gin/pkg/model/book"
)

type Book interface {
	GetByID(ctx context.Context, id int) (book_model.Book, error)
}
