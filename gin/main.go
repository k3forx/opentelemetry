package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k3forx/opentelemetry/gin/api/handler"
	"github.com/k3forx/opentelemetry/gin/opentelemetry/trace"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
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

	r := gin.New()
	r.Use(otelgin.Middleware("my-server"))

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		name := getUser(c, id)
		c.JSON(http.StatusOK, map[string]any{
			"id":   id,
			"name": name,
		})
	})

	v1 := r.Group("/v1")
	handler.RegisterBookHandler(v1)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func getUser(ctx *gin.Context, id string) string {
	// Pass the built-in `context.Context` object from http.Request to OpenTelemetry APIs
	// where required. It is available from gin.Context.Request.Context()
	_, span := trace.Tracer.Start(ctx.Request.Context(), "getUser", oteltrace.WithAttributes(attribute.String("id", id)))
	defer span.End()

	if id == "123" {
		return "tester"
	}
	return "unknown"
}
