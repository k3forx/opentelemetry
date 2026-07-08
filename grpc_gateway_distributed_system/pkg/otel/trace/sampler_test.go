package trace

import (
	"context"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// gateway の otelhttp が生成するルートスパンを模したパラメータで ShouldSample を直接呼ぶ
func serverSpanParams(path string) sdktrace.SamplingParameters {
	return sdktrace.SamplingParameters{
		ParentContext: context.Background(),
		Name:          "HTTP GET " + path,
		Kind:          oteltrace.SpanKindServer,
		Attributes:    []attribute.KeyValue{semconv.URLPathKey.String(path)},
	}
}

func TestDropPathSampler(t *testing.T) {
	s := &dropPathSampler{path: "/v1/health", next: sdktrace.AlwaysSample()}

	if got := s.ShouldSample(serverSpanParams("/v1/health")).Decision; got != sdktrace.Drop {
		t.Errorf("health span: got decision %v, want Drop", got)
	}
	if got := s.ShouldSample(serverSpanParams("/v1/books/1")).Decision; got != sdktrace.RecordAndSample {
		t.Errorf("books span: got decision %v, want RecordAndSample", got)
	}
}

// gateway でクライアントスパンを Drop すると、伝搬先の下流サービス（ParentBased）でも
// スパンがサンプリングされなくなることを、実際の TracerProvider と propagator で確認する
func TestDropGRPCClientSampler_DownstreamEffect(t *testing.T) {
	gatewayTP := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(&dropGRPCClientSampler{next: sdktrace.ParentBased(sdktrace.AlwaysSample())}),
	)
	downstreamTP := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)
	propagator := propagation.TraceContext{}

	// 1. gateway でルートスパン（サーバー）を生成 → サンプリングされる
	rootCtx, rootSpan := gatewayTP.Tracer("test").Start(context.Background(), "HTTP GET /v1/books/1",
		oteltrace.WithSpanKind(oteltrace.SpanKindServer))
	defer rootSpan.End()
	if !rootSpan.SpanContext().IsSampled() {
		t.Fatal("root span should be sampled")
	}

	// 2. gateway で gRPC クライアントスパンを生成 → Drop される
	clientCtx, clientSpan := gatewayTP.Tracer("test").Start(rootCtx, "book.v1.BookService/GetBook",
		oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	defer clientSpan.End()
	if clientSpan.SpanContext().IsSampled() {
		t.Fatal("client span should NOT be sampled")
	}

	// 3. クライアントスパンのコンテキストが W3C traceparent として伝搬される
	//    （otelgrpc が gRPC メタデータで行うことを propagator で直接再現）
	carrier := propagation.MapCarrier{}
	propagator.Inject(clientCtx, carrier)
	traceparent := carrier.Get("traceparent")
	if !strings.HasSuffix(traceparent, "-00") {
		t.Fatalf("traceparent should have sampled flag 00, got %q", traceparent)
	}

	// 4. 下流サービスで伝搬されたコンテキストからサーバースパンを生成
	//    → リモート親が未サンプリングなので ParentBased は Drop する
	downstreamCtx := propagator.Extract(context.Background(), carrier)
	_, serverSpan := downstreamTP.Tracer("test").Start(downstreamCtx, "GetBook",
		oteltrace.WithSpanKind(oteltrace.SpanKindServer))
	defer serverSpan.End()
	if serverSpan.SpanContext().IsSampled() {
		t.Fatal("downstream server span should NOT be sampled")
	}
}

// 対照実験: クライアントスパンを Drop しなければ下流もサンプリングされる
func TestParentBased_DownstreamSampled(t *testing.T) {
	gatewayTP := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)
	downstreamTP := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)
	propagator := propagation.TraceContext{}

	rootCtx, rootSpan := gatewayTP.Tracer("test").Start(context.Background(), "HTTP GET /v1/books/1",
		oteltrace.WithSpanKind(oteltrace.SpanKindServer))
	defer rootSpan.End()
	clientCtx, clientSpan := gatewayTP.Tracer("test").Start(rootCtx, "book.v1.BookService/GetBook",
		oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	defer clientSpan.End()
	if !clientSpan.SpanContext().IsSampled() {
		t.Fatal("client span should be sampled")
	}

	carrier := propagation.MapCarrier{}
	propagator.Inject(clientCtx, carrier)
	if !strings.HasSuffix(carrier.Get("traceparent"), "-01") {
		t.Fatalf("traceparent should have sampled flag 01, got %q", carrier.Get("traceparent"))
	}

	downstreamCtx := propagator.Extract(context.Background(), carrier)
	_, serverSpan := downstreamTP.Tracer("test").Start(downstreamCtx, "GetBook",
		oteltrace.WithSpanKind(oteltrace.SpanKindServer))
	defer serverSpan.End()
	if !serverSpan.SpanContext().IsSampled() {
		t.Fatal("downstream server span should be sampled")
	}
}

// 特定 API（/v1/authors 配下）だけ 100% サンプリングし、それ以外は
// TraceIDRatioBased(0) で全て Drop する構成の一気通貫の確認
func TestPriorityAPISampler(t *testing.T) {
	gatewayTP := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(&priorityAPISampler{
			pathPrefix: "/v1/authors",
			next:       sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0)),
		}),
	)
	downstreamTP := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0))),
	)
	propagator := propagation.TraceContext{}

	// 対象 API はルートスパンから下流まで全てサンプリングされる
	rootCtx, rootSpan := gatewayTP.Tracer("test").Start(context.Background(), "HTTP GET /v1/authors/1",
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(semconv.URLPathKey.String("/v1/authors/1")))
	defer rootSpan.End()
	if !rootSpan.SpanContext().IsSampled() {
		t.Fatal("authors root span should be sampled")
	}

	clientCtx, clientSpan := gatewayTP.Tracer("test").Start(rootCtx, "author.v1.AuthorService/GetAuthor",
		oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	defer clientSpan.End()
	if !clientSpan.SpanContext().IsSampled() {
		t.Fatal("authors client span should be sampled")
	}

	carrier := propagation.MapCarrier{}
	propagator.Inject(clientCtx, carrier)
	downstreamCtx := propagator.Extract(context.Background(), carrier)
	_, serverSpan := downstreamTP.Tracer("test").Start(downstreamCtx, "GetAuthor",
		oteltrace.WithSpanKind(oteltrace.SpanKindServer))
	defer serverSpan.End()
	if !serverSpan.SpanContext().IsSampled() {
		t.Fatal("downstream authors server span should be sampled")
	}

	// 対象外の API はルートスパンから Drop される
	booksCtx, booksSpan := gatewayTP.Tracer("test").Start(context.Background(), "HTTP GET /v1/books/1",
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(semconv.URLPathKey.String("/v1/books/1")))
	defer booksSpan.End()
	if booksSpan.SpanContext().IsSampled() {
		t.Fatal("books root span should NOT be sampled")
	}
	_ = booksCtx
}
