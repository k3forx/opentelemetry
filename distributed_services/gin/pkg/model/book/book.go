package book_model

type Book struct {
	ID       int64
	Title    string
	AuthorID int64
}
type BookWithAuthor struct {
	ID         int64
	Title      string
	AuthorID   int64
	AuthorName string
}
