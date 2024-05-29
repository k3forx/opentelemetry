package author_client

import (
	"context"
	"encoding/json"
	"fmt"

	author_model "github.com/k3forx/opentelemetry/gin/pkg/model/author"
)

type getAuthorByIDResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func (c *client) GetAuthorByID(ctx context.Context, id int64) (author_model.Author, error) {
	res, err := c.httpClient.Get(fmt.Sprintf("http://author-server:8081/authors/%d", id))
	if err != nil {
		return author_model.Author{}, fmt.Errorf("http get error: %w", err)
	}
	defer res.Body.Close()

	var getAuthorByIDResponse getAuthorByIDResponse
	if err := json.NewDecoder(res.Body).Decode(&getAuthorByIDResponse); err != nil {
		return author_model.Author{}, fmt.Errorf("decoding response error: %w", err)
	}
	return author_model.Author{
		ID:   getAuthorByIDResponse.ID,
		Name: getAuthorByIDResponse.Name,
		Bio:  getAuthorByIDResponse.Bio,
	}, nil
}
