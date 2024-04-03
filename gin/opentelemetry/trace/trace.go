package trace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewHTTPTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithInsecure(),
	}
	exp, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
	)
	defer func() {
		if err := traceProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	return traceProvider, nil
}
