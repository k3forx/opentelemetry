package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k3forx/opentelemetry/gin/api/handler"
	"github.com/k3forx/opentelemetry/gin/opentelemetry/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// func newExporter() (sdktrace.SpanExporter, error) {
// 	f, err := os.Create("trace.log")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return stdouttrace.New(
// 		stdouttrace.WithPrettyPrint(),
// 		stdouttrace.WithWriter(f),
// 	)
// }

// func newTraceProvider() (*sdktrace.TracerProvider, error) {
// 	exporter, err := newExporter()
// 	if err != nil {
// 		return nil, err
// 	}

// 	traceProvider := sdktrace.NewTracerProvider(
// 		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(time.Second)),
// 		sdktrace.WithSampler(sdktrace.AlwaysSample()),
// 	)
// 	otel.SetTracerProvider(sdktrace.NewTracerProvider())
// 	return traceProvider, nil
// }

func main() {
	ctx := context.Background()
	router := gin.Default()

	traceProvider, err := trace.NewHTTPTraceProvider(ctx)
	if err != nil {
		panic(err)
	}

	opts := otelgin.WithTracerProvider(traceProvider)
	router.Use(otelgin.Middleware("service-name", opts))

	v1 := router.Group("/v1")
	handler.RegisterBookHandler(v1)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
