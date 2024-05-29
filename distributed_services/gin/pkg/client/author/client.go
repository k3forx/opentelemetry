package author_client

import (
	"context"
	"net/http"

	author_model "github.com/k3forx/opentelemetry/gin/pkg/model/author"
)

type Client interface {
	GetAuthorByID(ctx context.Context, id int64) (author_model.Author, error)
}

type client struct {
	httpClient *http.Client
}

var _ Client = &client{}

func NewClient() *client {
	return &client{
		httpClient: &http.Client{},
	}
}
