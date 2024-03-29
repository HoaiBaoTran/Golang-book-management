package repository

import "github.com/hoaibao/book-management/pkg/domain"

type BookRepository interface {
	GetAllBooks(isbn, author, fromValue, toValue string) ([]domain.Book, error)
	GetBookById(id int) (domain.Book, error)
	CreateBook(book domain.Book, author []string) (domain.Book, error)
	DeleteBookById(bookId int) (domain.Book, error)
	UpdateBookById(bookId int, bookData map[string][]string) (domain.Book, error)
}
