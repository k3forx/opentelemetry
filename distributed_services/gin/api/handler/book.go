package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/k3forx/opentelemetry/gin/api/usecase"
	"github.com/k3forx/opentelemetry/gin/api/usecase/book/get_by_id"
	"github.com/k3forx/opentelemetry/gin/opentelemetry/trace"
	"github.com/k3forx/opentelemetry/gin/pkg/repository"

	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type bookHandler struct {
	repositorySet repository.RepositorySet
}

func newHandler(rs repository.RepositorySet) bookHandler {
	return bookHandler{
		repositorySet: rs,
	}
}

func RegisterBookHandler(group *gin.RouterGroup, rs repository.RepositorySet) {
	h := newHandler(rs)
	group.GET("/books/:id", h.GetByID)
}

func (h bookHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	ctx, span := trace.Tracer.Start(
		c.Request.Context(), trace.SpanNameHandler,
		oteltrace.WithAttributes(
			attribute.String("name", "GetByID"),
			attribute.String("id", idStr),
		),
	)
	defer span.End()

	u := get_by_id.NewUsecase(h.repositorySet)
	id, _ := strconv.Atoi(idStr)

	executer := usecase.NewUsecaseExecuter(u)
	out := executer.DoWithTrace(ctx, get_by_id.Input{ID: int64(id)})
	if out.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": out.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         out.Book.ID,
		"title":      out.Book.Title,
		"authorName": out.Author.Name,
	})
}
