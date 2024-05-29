package repository

import (
	"database/sql"

	author_client "github.com/k3forx/opentelemetry/gin/pkg/client/author"
	author_repository "github.com/k3forx/opentelemetry/gin/pkg/repository/author"
	book_repository "github.com/k3forx/opentelemetry/gin/pkg/repository/book"
	author_repository_impl "github.com/k3forx/opentelemetry/gin/pkg/repositoryimpl/author"
	book_repository_impl "github.com/k3forx/opentelemetry/gin/pkg/repositoryimpl/book"
)

type RepositorySet struct {
	Author       author_repository.AuthorRepository
	AuthorClient author_client.Client
	Book         book_repository.BookRepository
}

func SetUp(db *sql.DB) RepositorySet {
	return RepositorySet{
		Author:       author_repository_impl.NewAuthorRepository(db),
		AuthorClient: author_client.NewClient(),
		Book:         book_repository_impl.NewBookRepository(db),
	}
}
