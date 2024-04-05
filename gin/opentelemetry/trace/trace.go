package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewHTTPTraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// opts := []otlptracehttp.Option{
	// 	otlptracehttp.WithInsecure(),
	// 	otlptracehttp.WithEndpointURL("http://0.0.0.0:4318/v1/traces"),
	// }
	// exp, err := otlptracehttp.New(ctx, opts...)
	// if err != nil {
	// 	return nil, err
	// }

	debugExp, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		// stdouttrace.WithWriter(os.Stderr),
	)
	if err != nil {
		return nil, err
	}

	traceProvider := sdktrace.NewTracerProvider(
		// sdktrace.WithBatcher(exp),
		sdktrace.WithBatcher(debugExp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return traceProvider, nil
}
