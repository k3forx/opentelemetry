package create

import (
	"context"
	"errors"
	"log"

	usecase_iface "github.com/k3forx/opentelemetry/gin/api/usecase"
	book_model "github.com/k3forx/opentelemetry/gin/pkg/model/book"
	"github.com/k3forx/opentelemetry/gin/pkg/repository"
	author_repository "github.com/k3forx/opentelemetry/gin/pkg/repository/author"
	book_repository "github.com/k3forx/opentelemetry/gin/pkg/repository/book"
)

type Input struct {
	AuthorID int64
	Title    string
}

type Output struct {
	Error error
	Book  book_model.Book
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
	author, err := u.authorRepository.GetByID(ctx, in.AuthorID)
	if err != nil {
		log.Panicln(err)
		return Output{Error: err}
	}
	if author.ID <= 0 {
		return Output{Error: errors.New("author is not found")}
	}

	book := book_model.Book{
		Title:    in.Title,
		AuthorID: author.ID,
	}
	if err := u.bookRepository.Create(ctx, &book); err != nil {
		return Output{Error: err}
	}

	return Output{Book: book}
}
