package repository

import "github.com/hoaibao/book-management/pkg/domain"

type BookRepository interface {
	GetAllBooksVersion2(fromValue, toValue string) ([]domain.Book, error)
	GetBookByIdVersion2(id int) (domain.Book, error)
	CreateBookVersion2(book domain.Book) (domain.Book, error)
	DeleteBookByIdVersion2(bookId int) (domain.Book, error)
	UpdateBookByIdVersion2(bookId int, bookData map[string]string) (domain.Book, error)
}
