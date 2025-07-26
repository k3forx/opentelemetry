package server

import (
	"context"

	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/api/usecase"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/api/usecase/book/get_by_id"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/otel/trace"
	bookv1 "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/proto/gen/book/v1"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookServer struct {
	bookv1.UnimplementedBookServiceServer
	repositorySet repository.RepositorySet
}

func NewBookServer(rs repository.RepositorySet) *BookServer {
	return &BookServer{
		repositorySet: rs,
	}
}

func (s *BookServer) GetBook(ctx context.Context, req *bookv1.GetBookRequest) (*bookv1.GetBookResponse, error) {
	ctx, span := trace.Tracer.Start(
		ctx, trace.SpanNameHandler,
		oteltrace.WithAttributes(
			attribute.String("name", "GetBook"),
			attribute.Int64("id", req.Id),
		),
	)
	defer span.End()

	u := get_by_id.NewUsecase(s.repositorySet)
	executer := usecase.NewUsecaseExecuter(u)
	out := executer.DoWithTrace(ctx, get_by_id.Input{ID: req.Id})
	if out.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get book: %v", out.Error)
	}

	return &bookv1.GetBookResponse{
		Book: &bookv1.Book{
			Id:         out.Book.ID,
			Title:      out.Book.Title,
			AuthorName: out.Author.Name,
		},
	}, nil
}
