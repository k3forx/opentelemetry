package main

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/pkg/otel/trace"
	authorv1 "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/proto/gen/author/v1"
	bookv1 "github.com/k3forx/opentelemetry/grpc_gateway_distributed_system/proto/gen/book/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// HealthResponse はヘルスチェックのレスポンス構造体
type HealthResponse struct {
	Status string `json:"status"`
}

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

// healthHandler はヘルスチェックのハンドラー
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{Status: "healthy"}
	json.NewEncoder(w).Encode(response)
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

	grpcMux := runtime.NewServeMux()

	if err := bookv1.RegisterBookServiceHandler(ctx, grpcMux, bookConn); err != nil {
		return fmt.Errorf("failed to register book service: %v", err)
	}

	if err := authorv1.RegisterAuthorServiceHandler(ctx, grpcMux, authorConn); err != nil {
		return fmt.Errorf("failed to register author service: %v", err)
	}

	const healthPath = "/v1/health"
	mux := http.NewServeMux()
	mux.HandleFunc(healthPath, healthHandler)
	mux.Handle("/", grpcMux)

	port := cmp.Or(os.Getenv("GATEWAY_PORT"), "8080")

	log.Printf("gateway server listening on :%s", port)

	handler := otelhttp.NewHandler(mux, "grpc-gateway",
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != healthPath
		}),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path)
		}),
		otelhttp.WithPublicEndpoint(), // セキュリティとベストプラクティスのため
	)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}
