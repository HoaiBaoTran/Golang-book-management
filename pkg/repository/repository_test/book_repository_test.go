package repository_test

import (
	"testing"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/hoaibao/book-management/pkg/domain"
	"github.com/hoaibao/book-management/pkg/repository"
)

func TestBookRepository_GetAllBooks(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	testBookRepository := repository.NewTestMemoryBookRepository(db)

	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "author", "publish_year"}).
		AddRow(1, "Book 1", "1234567890", "Author 1", "2011").
		AddRow(2, "Book 2", "0987654321", "Author 2", "2012")

	mock.ExpectQuery("^SELECT \\* FROM book$").
		WillReturnRows(rows)

	_, err = testBookRepository.GetAllBooks("", "", "", "")
	if err != nil {
		t.Errorf("Error creating book: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBookRepository_GetAllBooksByFilter(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	testBookRepository := repository.NewTestMemoryBookRepository(db)

	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "author", "publish_year"}).
		AddRow(1, "Book 1", "1234567890", "Author 1", "2011").
		AddRow(2, "Book 2", "0987654321", "Author 2", "2013")

	mock.ExpectQuery(`^SELECT \* FROM book WHERE isbn = '1234567890' AND author = 'Author 1' AND publish_year >= 2011 AND publish_year <= 2012`).
		WillReturnRows(rows)

	_, err = testBookRepository.GetAllBooks("1234567890", "Author 1", "2011", "2012")
	if err != nil {
		t.Errorf("Error creating book: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBookRepository_GetBookById(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	testBookRepository := repository.NewTestMemoryBookRepository(db)

	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "author", "publish_year"}).
		AddRow(1, "Book 1", "1234567890", "Author 1", "2011").
		AddRow(2, "Book 2", "0987654321", "Author 2", "2012")

	mock.ExpectQuery("^SELECT \\* FROM book WHERE id = \\$1$").
		WithArgs(1).
		WillReturnRows(rows)

	_, err = testBookRepository.GetBookById(1)
	if err != nil {
		t.Errorf("Error creating book: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBookRepository_DeleteBookById(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	testBookRepository := repository.NewTestMemoryBookRepository(db)

	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "author", "publish_year"}).
		AddRow(1, "Book 1", "1234567890", "Author 1", "2011").
		AddRow(2, "Book 2", "0987654321", "Author 2", "2012")

	mock.ExpectQuery("^SELECT \\* FROM book WHERE id = \\$1$").
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectExec("^DELETE FROM book WHERE id = \\$1$").
		WithArgs(1).
		WillReturnResult(sqlMock.NewResult(0, 1))

	_, err = testBookRepository.DeleteBookById(1)
	if err != nil {
		t.Errorf("Error creating book: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBookRepository_UpdateBookById(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	testBookRepository := repository.NewTestMemoryBookRepository(db)

	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "author", "publish_year"}).
		AddRow(1, "Book 1", "1234567890", "Author 1", "2011").
		AddRow(2, "Book 2", "0987654321", "Author 2", "2012")

	bookData := map[string]string{
		"name":         "Updated book",
		"isbn":         "0987654321",
		"author":       "Updated author",
		"publish_year": "2014",
	}

	mock.ExpectQuery("^SELECT \\* FROM book WHERE id = \\$1$").
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectExec("^UPDATE book SET name = 'Updated book', isbn = '0987654321', author = 'Updated author', publish_year = '2014' WHERE id = \\$1$").
		WithArgs(1).
		WillReturnResult(sqlMock.NewResult(1, 1))

	_, err = testBookRepository.UpdateBookById(1, bookData)
	if err != nil {
		t.Errorf("Error creating book: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestBookRepository_CreateBook(t *testing.T) {

	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	book := domain.Book{Id: 1, Name: "Capybara", ISBN: "1234567890", Author: "Capybara Writer", PublishYear: 2024}
	mock.ExpectExec("^INSERT INTO book\\(isbn, name, author, publish_year\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)$").
		WithArgs(book.ISBN, book.Name, book.Author, book.PublishYear).
		WillReturnResult(sqlMock.NewResult(1, 1))

	_, err = testBookRepository.CreateBook(book)
	if err != nil {
		t.Errorf("Error creating book: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
