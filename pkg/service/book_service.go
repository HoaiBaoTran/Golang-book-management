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

func (s *BookService) GetAllBooksVersion2(fromValue, toValue string) ([]domain.Book, error) {
	return s.bookRepository.GetAllBooksVersion2(fromValue, toValue)
}

func (s *BookService) GetBookByIdVersion2(id int) (domain.Book, error) {
	return s.bookRepository.GetBookByIdVersion2(id)
}

func (s *BookService) CreateBookVersion2(name, isbn, author string, publishYear int) (domain.Book, error) {
	book := domain.Book{
		ISBN: isbn,
		Name: name,
		Author: domain.Author{
			Name: author,
		},
		PublishYear: publishYear,
	}

	return s.bookRepository.CreateBookVersion2(book)
}

func (s *BookService) DeleteBookByIdVersion2(bookId int) (domain.Book, error) {
	return s.bookRepository.DeleteBookByIdVersion2(bookId)
}

func (s *BookService) UpdateBookByIdVersion2(bookId int, bookData map[string]string) (domain.Book, error) {
	return s.bookRepository.UpdateBookByIdVersion2(bookId, bookData)
}
