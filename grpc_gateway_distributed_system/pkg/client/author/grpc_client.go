package author_client

import (
	"context"
	"fmt"

	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/opentelemetry/trace"
	author_model "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/model/author"
	authorv1 "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/proto/gen/author/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	GetAuthorByID(ctx context.Context, id int64) (author_model.Author, error)
}

type grpcClient struct {
	client authorv1.AuthorServiceClient
	conn   *grpc.ClientConn
}

var _ Client = &grpcClient{}

func NewGRPCClient(grpcServerAddr string) (Client, error) {
	conn, err := grpc.NewClient(
		grpcServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC server: %w", err)
	}

	client := authorv1.NewAuthorServiceClient(conn)
	return &grpcClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *grpcClient) GetAuthorByID(ctx context.Context, id int64) (author_model.Author, error) {
	ctx, span := trace.Tracer.Start(
		ctx, trace.SpanNameHandler,
		oteltrace.WithAttributes(
			attribute.String("name", "GetAuthorByID"),
			attribute.Int64("id", id),
		),
	)
	defer span.End()

	req := &authorv1.GetAuthorRequest{
		Id: id,
	}

	resp, err := c.client.GetAuthor(ctx, req)
	if err != nil {
		return author_model.Author{}, fmt.Errorf("failed to get author via gRPC: %w", err)
	}

	return author_model.Author{
		ID:   resp.Author.Id,
		Name: resp.Author.Name,
		Bio:  resp.Author.Bio,
	}, nil
}

func (c *grpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
