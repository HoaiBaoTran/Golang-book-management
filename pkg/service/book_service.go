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

func (s *BookService) CreateBook(name, isbn string, author []string, publishYear int) (domain.Book, error) {
	authorObj := []domain.Author{}

	for _, authorName := range author {
		authorObj = append(authorObj, domain.Author{Name: authorName})
	}

	book := domain.Book{
		ISBN:        isbn,
		Name:        name,
		Authors:     authorObj,
		PublishYear: publishYear,
	}

	return s.bookRepository.CreateBook(book, author)
}

func (s *BookService) DeleteBookById(bookId int) (domain.Book, error) {
	return s.bookRepository.DeleteBookById(bookId)
}

func (s *BookService) UpdateBookById(bookId int, bookData map[string][]string) (domain.Book, error) {
	return s.bookRepository.UpdateBookById(bookId, bookData)
}
