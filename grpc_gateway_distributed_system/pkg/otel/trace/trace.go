package trace

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
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
		sdktrace.WithSampler(newSampler()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp, nil
}

type sampler struct {
	defaultSampler sdktrace.Sampler
}

func newSampler() sdktrace.Sampler {
	return &sampler{
		defaultSampler: sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.01)),
	}
}

func (s *sampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	spanContext := trace.SpanContextFromContext(p.ParentContext)
	fmt.Println("--------------------------------")
	fmt.Printf("Name:          %s\n", p.Name)
	fmt.Printf("ParentContext: %+v\n", p.ParentContext)

	for _, attr := range p.Attributes {
		fmt.Printf("Key: %s, Value: %s\n", attr.Key, attr.Value.AsString())
		if attr.Key == semconv.HTTPRouteKey {
			if attr.Value.AsString() == "/authors/:id" {
				fmt.Println("Sampling!!!!!!!")
				return sdktrace.SamplingResult{
					Decision:   sdktrace.RecordAndSample,
					Attributes: p.Attributes,
					Tracestate: spanContext.TraceState(),
				}
			}
		}
		// if attr.Key == "name" {
		// 	return sdktrace.SamplingResult{
		// 		Decision:   sdktrace.RecordAndSample,
		// 		Attributes: p.Attributes,
		// 		Tracestate: spanContext.TraceState(),
		// 	}
		// }
	}
	fmt.Println("--------------------------------")

	return s.defaultSampler.ShouldSample(p)
}

func (s *sampler) Description() string {
	return s.defaultSampler.Description()
}
