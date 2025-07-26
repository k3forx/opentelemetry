package get_by_id

import (
	"context"

	usecase_iface "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/api/usecase"
	author_model "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/model/author"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository"
	author_repository "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository/author"
)

type Input struct {
	ID int64
}

type Output struct {
	Error  error
	Author author_model.Author
}

type usecase struct {
	authorRepository author_repository.AuthorRepository
}

func NewUsecase(repositorySet repository.RepositorySet) usecase_iface.Usecase[Input, Output] {
	return usecase{
		authorRepository: repositorySet.Author,
	}
}

func (u usecase) Do(ctx context.Context, in Input) Output {
	author, err := u.authorRepository.GetByID(ctx, in.ID)
	if err != nil {
		return Output{Error: err}
	}
	return Output{
		Author: author_model.Author{
			ID:   author.ID,
			Name: author.Name,
		},
	}
}
