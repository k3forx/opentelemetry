package book_repository

import (
	"context"

	book_model "github.com/k3forx/opentelemetry/gin/pkg/model/book"
)

type BookRepository interface {
	Create(ctx context.Context, book *book_model.Book) error
	GetAllByAuthorID(ctx context.Context, id int64) ([]book_model.Book, error)
	GetWithAuthorByID(ctx context.Context, id int64) (book_model.BookWithAuthor, error)
	GetByID(ctx context.Context, id int64) (book_model.Book, error)
}
