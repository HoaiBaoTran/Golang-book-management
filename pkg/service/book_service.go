package service

import (
	"github.com/hoaibao/book-management/pkg/domain"
	"github.com/hoaibao/book-management/pkg/repository"
)

type BookService struct {
	bookRepository repository.BookRepository
}

func NewBookService(bookRepository repository.BookRepository) *BookService {
	return &BookService{
		bookRepository: bookRepository,
	}
}

func (s *BookService) GetAllBooks(isbn, author, fromValue, toValue string) ([]domain.Book, error) {
	return s.bookRepository.GetAllBooks(isbn, author, fromValue, toValue)
}

func (s *BookService) GetBookById(id int) (domain.Book, error) {
	return s.bookRepository.GetBookById(id)
}

func (s *BookService) CreateBook(name, isbn, author string, publishYear int) (domain.Book, error) {
	book := domain.Book{
		ISBN:        isbn,
		Name:        name,
		Author:      author,
		PublishYear: publishYear,
	}

	return s.bookRepository.CreateBook(book)
}

func (s *BookService) DeleteBookById(bookId int) (domain.Book, error) {
	return s.bookRepository.DeleteBookById(bookId)
}

func (s *BookService) UpdateBookById(bookId int, bookData map[string]string) (domain.Book, error) {
	return s.bookRepository.UpdateBookById(bookId, bookData)
}
