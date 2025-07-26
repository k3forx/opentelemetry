package get_by_id

import (
	"context"

	usecase_iface "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/api/usecase"
	author_client "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/client/author"
	author_model "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/model/author"
	book_model "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/model/book"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository"
	book_repository "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository/book"
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
	bookRepository book_repository.BookRepository
	authorClient   author_client.Client
}

func NewUsecase(repositorySet repository.RepositorySet) usecase_iface.Usecase[Input, Output] {
	return usecase{
		bookRepository: repositorySet.Book,
		authorClient:   repositorySet.AuthorClient,
	}
}

func (u usecase) Do(ctx context.Context, in Input) Output {
	book, err := u.bookRepository.GetByID(ctx, in.ID)
	if err != nil {
		return Output{Error: err}
	}

	author, err := u.authorClient.GetAuthorByID(ctx, book.AuthorID)
	if err != nil {
		return Output{Error: err}
	}

	return Output{
		Book: book_model.Book{
			ID:    book.ID,
			Title: book.Title,
		},
		Author: author_model.Author{
			ID:   book.AuthorID,
			Name: author.Name,
		},
	}
}
