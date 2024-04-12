package get_by_id

import (
	"context"

	usecase_iface "github.com/k3forx/opentelemetry/gin/api/usecase"
	author_model "github.com/k3forx/opentelemetry/gin/pkg/model/author"
	book_model "github.com/k3forx/opentelemetry/gin/pkg/model/book"
	"github.com/k3forx/opentelemetry/gin/pkg/repository"
	author_repository "github.com/k3forx/opentelemetry/gin/pkg/repository/author"
	book_repository "github.com/k3forx/opentelemetry/gin/pkg/repository/book"
)

type Input struct {
	ID int64
}

type Output struct {
	Error  error
	Book   book_model.Book
	Author author_model.Author
}

type usecase struct {
	authorRepository author_repository.AuthorRepository
	bookRepository   book_repository.BookRepository
}

func NewUsecase(repositorySet repository.RepositorySet) usecase_iface.Usecase[Input, Output] {
	return usecase{
		authorRepository: repositorySet.Author,
		bookRepository:   repositorySet.Book,
	}
}

func (u usecase) Do(ctx context.Context, in Input) Output {
	book, err := u.bookRepository.GetByID(ctx, in.ID)
	if err != nil {
		return Output{Error: err}
	}

	author, err := u.authorRepository.GetByID(ctx, book.AuthorID)
	if err != nil {
		return Output{Error: err}
	}

	return Output{
		Book:   book,
		Author: author,
	}
}
