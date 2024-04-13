package book_repository_impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/k3forx/opentelemetry/gin/opentelemetry/trace"
	book_model "github.com/k3forx/opentelemetry/gin/pkg/model/book"
	book_repository "github.com/k3forx/opentelemetry/gin/pkg/repository/book"
)

var _ book_repository.BookRepository = &bookRepositoryImpl{}

type bookRepositoryImpl struct {
	queries *Queries
}

func NewBookRepository(db *sql.DB) *bookRepositoryImpl {
	return &bookRepositoryImpl{queries: New(db)}
}

func (impl *bookRepositoryImpl) Create(ctx context.Context, book *book_model.Book) error {
	ctx, span := trace.Tracer.Start(ctx, trace.SpanNameRepository)
	defer span.End()

	res, err := impl.queries.Create(ctx, CreateParams{
		AuthorID: book.AuthorID,
		Title:    book.Title,
	})
	if err != nil {
		return fmt.Errorf("create book error: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id error: %w", err)
	}
	book.ID = id

	return nil
}

func (impl *bookRepositoryImpl) GetAllByAuthorID(ctx context.Context, id int64) ([]book_model.Book, error) {
	ctx, span := trace.Tracer.Start(ctx, trace.SpanNameRepository)
	defer span.End()

	bs, err := impl.queries.GetAllByAuthorID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []book_model.Book{}, nil
		}
		return nil, err
	}

	ms := make([]book_model.Book, len(bs))
	for i, b := range bs {
		ms[i] = book_model.Book{
			ID:       b.ID,
			Title:    b.Title,
			AuthorID: b.AuthorID,
		}
	}
	return ms, nil
}

func (impl *bookRepositoryImpl) GetWithAuthorByID(ctx context.Context, id int64) (book_model.BookWithAuthor, error) {
	ctx, span := trace.Tracer.Start(ctx, trace.SpanNameRepository)
	defer span.End()

	res, err := impl.queries.GetWithAuthorByID(ctx, id)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return book_model.BookWithAuthor{}, nil
		}
		return book_model.BookWithAuthor{}, err
	}

	return book_model.BookWithAuthor{
		ID:         res.ID,
		Title:      res.Title,
		AuthorID:   res.AuthorID,
		AuthorName: res.AuthorName.String,
	}, nil
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
