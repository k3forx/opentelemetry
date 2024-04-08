package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	usecase_iface "github.com/k3forx/opentelemetry/gin/api/usecase"
	"github.com/k3forx/opentelemetry/gin/api/usecase/book/get_by_id"
	"github.com/k3forx/opentelemetry/gin/opentelemetry/trace"

	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type bookHandler struct{}

func newHandler() bookHandler {
	return bookHandler{}
}

func RegisterBookHandler(group *gin.RouterGroup) {
	h := newHandler()
	group.GET("/books/:id", h.GetByID)
}

func (h bookHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	ctx, span := trace.Tracer.Start(
		c.Request.Context(), "GetByID",
		oteltrace.WithAttributes(attribute.String("id", idStr)),
	)
	defer span.End()

	u := get_by_id.NewUsecase()
	id, _ := strconv.Atoi(idStr)

	executer := usecase_iface.NewUsecaseExecuter(u)
	out := executer.DoWithTrace(ctx, get_by_id.Input{ID: id})

	c.JSON(http.StatusOK, gin.H{
		"id":   out.ID,
		"name": out.Name,
	})
}
