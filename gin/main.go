package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/k3forx/opentelemetry/gin/api/handler"
	"github.com/k3forx/opentelemetry/gin/opentelemetry/trace"
	"github.com/k3forx/opentelemetry/gin/pkg/repository"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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

	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root_password",
		DBName:               "app",
		Addr:                 "mysql:3306",
		Net:                  "tcp",
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	repositorySet := repository.SetUp(db)

	r := gin.New()
	r.Use(otelgin.Middleware("my-server"))

	v1 := r.Group("/v1")
	handler.RegisterBookHandler(v1, repositorySet)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
