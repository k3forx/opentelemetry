package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type bookHandler struct{}

func newHandler() bookHandler {
	return bookHandler{}
}

func RegisterBookHandler(group *gin.RouterGroup) {
	h := newHandler()
	group.GET("/books/:id", h.GetByID)
}

func (h bookHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	ctx.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
