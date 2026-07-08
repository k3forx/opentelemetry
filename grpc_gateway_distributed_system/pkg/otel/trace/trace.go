package trace

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var Tracer = otel.Tracer("grpc-gateway-distributed-system")

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

// newSampler は OTEL_SAMPLER 環境変数に応じたサンプラーを返す（検証用）
func newSampler() sdktrace.Sampler {
	var s sdktrace.Sampler
	switch os.Getenv("OTEL_SAMPLER") {
	case "drop_health":
		s = &dropPathSampler{path: "/v1/health", next: sdktrace.AlwaysSample()}
	case "drop_grpc_client":
		s = &dropGRPCClientSampler{next: sdktrace.ParentBased(sdktrace.AlwaysSample())}
	case "priority_api":
		s = &priorityAPISampler{pathPrefix: "/v1/authors", next: sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0))}
	case "parent_based":
		s = sdktrace.ParentBased(sdktrace.AlwaysSample())
	default:
		s = sdktrace.AlwaysSample()
	}
	return &loggingSampler{next: s}
}

// dropPathSampler は url.path が path に一致するサーバースパンを Drop し、
// それ以外は next に委譲するサンプラー
type dropPathSampler struct {
	path string
	next sdktrace.Sampler
}

func (s *dropPathSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if p.Kind == oteltrace.SpanKindServer {
		for _, attr := range p.Attributes {
			if attr.Key == semconv.URLPathKey && attr.Value.AsString() == s.path {
				return sdktrace.SamplingResult{
					Decision:   sdktrace.Drop,
					Tracestate: oteltrace.SpanContextFromContext(p.ParentContext).TraceState(),
				}
			}
		}
	}
	return s.next.ShouldSample(p)
}

func (s *dropPathSampler) Description() string {
	return fmt.Sprintf("DropPathSampler{path:%s}", s.path)
}

// dropGRPCClientSampler はクライアントスパンを Drop し、それ以外は next に委譲するサンプラー
type dropGRPCClientSampler struct {
	next sdktrace.Sampler
}

func (s *dropGRPCClientSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if p.Kind == oteltrace.SpanKindClient {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.Drop,
			Tracestate: oteltrace.SpanContextFromContext(p.ParentContext).TraceState(),
		}
	}
	return s.next.ShouldSample(p)
}

func (s *dropGRPCClientSampler) Description() string {
	return "DropGRPCClientSampler"
}

// priorityAPISampler は url.path が pathPrefix から始まるルートスパンを常にサンプリングし、
// それ以外は next に委譲するサンプラー
type priorityAPISampler struct {
	pathPrefix string
	next       sdktrace.Sampler
}

func (s *priorityAPISampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	psc := oteltrace.SpanContextFromContext(p.ParentContext)
	// ルートスパン（親スパンなし）のみ独自に判定する
	if p.Kind == oteltrace.SpanKindServer && !psc.IsValid() {
		for _, attr := range p.Attributes {
			if attr.Key == semconv.URLPathKey && strings.HasPrefix(attr.Value.AsString(), s.pathPrefix) {
				return sdktrace.SamplingResult{
					Decision:   sdktrace.RecordAndSample,
					Tracestate: psc.TraceState(),
				}
			}
		}
	}
	return s.next.ShouldSample(p)
}

func (s *priorityAPISampler) Description() string {
	return fmt.Sprintf("PriorityAPISampler{pathPrefix:%s}", s.pathPrefix)
}

// loggingSampler は ShouldSample に渡ってくる情報と判定結果を観察するためのデバッグ用サンプラー
type loggingSampler struct {
	next sdktrace.Sampler
}

func (s *loggingSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	res := s.next.ShouldSample(p)
	psc := oteltrace.SpanContextFromContext(p.ParentContext)
	fmt.Printf("[sampler] name=%q kind=%s parent(valid=%t remote=%t sampled=%t) decision=%d\n",
		p.Name, p.Kind, psc.IsValid(), psc.IsRemote(), psc.IsSampled(), res.Decision)
	for _, attr := range p.Attributes {
		fmt.Printf("[sampler]   attr %s=%v\n", attr.Key, attr.Value.Emit())
	}
	return res
}

func (s *loggingSampler) Description() string {
	return s.next.Description()
}
