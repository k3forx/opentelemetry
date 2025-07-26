package server

import (
	"context"

	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/api/usecase"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/api/usecase/author/get_by_id"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/opentelemetry/trace"
	authorv1 "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/proto/gen/author/v1"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/repository"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorServer struct {
	authorv1.UnimplementedAuthorServiceServer
	repositorySet repository.RepositorySet
}

func NewAuthorServer(rs repository.RepositorySet) *AuthorServer {
	return &AuthorServer{
		repositorySet: rs,
	}
}

func (s *AuthorServer) GetAuthor(ctx context.Context, req *authorv1.GetAuthorRequest) (*authorv1.GetAuthorResponse, error) {
	ctx, span := trace.Tracer.Start(
		ctx, trace.SpanNameHandler,
		oteltrace.WithAttributes(
			attribute.String("name", "GetAuthor"),
			attribute.Int64("id", req.Id),
		),
	)
	defer span.End()

	u := get_by_id.NewUsecase(s.repositorySet)
	executer := usecase.NewUsecaseExecuter(u)
	out := executer.DoWithTrace(ctx, get_by_id.Input{ID: req.Id})
	if out.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get author: %v", out.Error)
	}

	return &authorv1.GetAuthorResponse{
		Author: &authorv1.Author{
			Id:   out.Author.ID,
			Name: out.Author.Name,
			Bio:  out.Author.Bio,
		},
	}, nil
}
