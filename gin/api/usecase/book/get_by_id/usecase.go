package get_by_id

import (
	"context"

	usecase_iface "github.com/k3forx/opentelemetry/gin/api/usecase"
	book_model "github.com/k3forx/opentelemetry/gin/pkg/model/book"
	book_repository "github.com/k3forx/opentelemetry/gin/pkg/repository/book"
)

type Input struct {
	ID int
}

type Output struct {
	ID   int
	Name string
}

type usecase struct {
	bookRepository book_repository.Book
}

func NewUsecase() usecase_iface.Usecase[Input, Output] {
	return usecase{}
}

func (u usecase) Do(ctx context.Context, in Input) Output {
	// book, err := u.bookRepository.GetByID(ctx, in.ID)
	// if err != nil {
	// 	return Output{}
	// }
	book := book_model.Book{
		ID:   in.ID,
		Name: "hogehoge",
	}

	return Output{
		ID:   book.ID,
		Name: book.Name,
	}
}
