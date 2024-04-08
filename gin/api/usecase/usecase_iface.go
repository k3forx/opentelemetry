package usecase_iface

import (
	"context"

	"github.com/k3forx/opentelemetry/gin/opentelemetry/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Usecase[T, V any] interface {
	Do(context.Context, T) V
}

type UsecaseExecuter[T, V any] interface {
	Do(ctx context.Context, T any) V
	DoWithTrace(ctx context.Context, T any, opts ...oteltrace.SpanStartOption) V
}

type executer[T, V any] struct {
	u Usecase[T, V]
}

func NewUsecaseExecuter[T, V any, u Usecase[T, V]](uc u) executer[T, V] {
	return executer[T, V]{
		u: uc,
	}
}

func (e *executer[T, V]) Do(ctx context.Context, in T) V {
	return e.u.Do(ctx, in)
}

func (e *executer[T, V]) DoWithTrace(ctx context.Context, in T, opts ...oteltrace.SpanStartOption) V {
	ctx, span := trace.Tracer.Start(ctx, "Usecase.Do", opts...)
	defer span.End()
	return e.Do(ctx, in)
}
