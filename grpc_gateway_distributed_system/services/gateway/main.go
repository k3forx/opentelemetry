package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/otel/trace"
	authorv1 "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/proto/gen/author/v1"
	bookv1 "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/proto/gen/book/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	tp, err := trace.InitTraceProvider(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// gRPC Gateway サーバーを起動
	if err := startGatewayServer(ctx); err != nil {
		log.Fatal("Failed to start gateway server:", err)
	}
}

func startGatewayServer(ctx context.Context) error {
	// Book Service への接続
	bookConn, err := grpc.NewClient(
		"book-service:9091",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial book service: %v", err)
	}
	defer bookConn.Close()

	// Author Service への接続
	authorConn, err := grpc.NewClient(
		"author-service:9092",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial author service: %v", err)
	}
	defer authorConn.Close()

	mux := runtime.NewServeMux()

	// Book Service の登録
	if err := bookv1.RegisterBookServiceHandler(ctx, mux, bookConn); err != nil {
		return fmt.Errorf("failed to register book service: %v", err)
	}

	// Author Service の登録
	if err := authorv1.RegisterAuthorServiceHandler(ctx, mux, authorConn); err != nil {
		return fmt.Errorf("failed to register author service: %v", err)
	}

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway server listening on :%s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
