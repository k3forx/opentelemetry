package trace

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var Tracer = otel.Tracer("gin-server")

const (
	SpanNameHandler    = "handler"
	SpanNameUsecase    = "usecase"
	SpanNameRepository = "repository"
)

func InitTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	traceHTTPOpts := []otlptracehttp.Option{
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpointURL("http://jaeger:4318/v1/traces"),
	}
	exp, err := otlptracehttp.New(ctx, traceHTTPOpts...)
	if err != nil {
		return nil, err
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(os.Getenv("OTEL_SERVICE_NAME")),
		semconv.ServiceVersionKey.String("0.0.1"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}
