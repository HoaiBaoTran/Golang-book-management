package repository_test

import (
	"fmt"
	"testing"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
	// "github.com/hoaibao/book-management/pkg/domain"

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

	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(1, "Book 1", "1234567890", "2011", "1", "Capybara", "01/01/2002").
		AddRow(2, "Book 2", "0987654321", "2012", "2", "Luffy", "02/02/2002")

	mock.ExpectQuery("^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id$").
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

	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(1, "Book 1", "1234567890", "2011", "1", "Capybara", "01/01/2002").
		AddRow(2, "Book 2", "0987654321", "2012", "2", "Luffy", "02/02/2002")

	mock.ExpectQuery(`^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id WHERE b.isbn = '1234567890' AND a.name = 'Author 1' AND b.publish_year >= 2011 AND b.publish_year <= 2012`).
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

	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(1, "Book 1", "1234567890", "2011", "1", "Capybara", "01/01/2002").
		AddRow(2, "Book 2", "0987654321", "2012", "2", "Luffy", "02/02/2002")

	mock.ExpectQuery("^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id WHERE b.id = \\$1$").
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
	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(1, "Book 1", "1234567890", "2011", "1", "Capybara", "01/01/2002")

	mock.ExpectQuery("^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id WHERE b.id = \\$1$").
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM book_author where book_id = \\$1$").
		WithArgs(1).
		WillReturnResult(sqlMock.NewResult(1, 1))

	mock.ExpectExec("^DELETE FROM book WHERE id = \\$1$").
		WithArgs(1).
		WillReturnResult(sqlMock.NewResult(1, 1))

	mock.ExpectCommit()

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

	rows := sqlMock.NewRows([]string{"b.id", "b.isbn", "b.name", "b.publish_year", "a.id", "a.name", "a.birth_day"}).
		AddRow(2, "0987654321", "Book 2", "2012", 1, "Capybara", "02/02/2002")

	bookData := map[string][]string{
		"name": {"Updated book"},
	}
	mock.ExpectQuery("^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id WHERE b.id = \\$1$").
		WithArgs(2).
		WillReturnRows(rows)

	mock.ExpectExec("^UPDATE book SET name = 'Updated book' WHERE id = \\$1").
		WithArgs(2).
		WillReturnResult(sqlMock.NewResult(1, 1))

	_, err = testBookRepository.UpdateBookById(2, bookData)
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

	sqlMock.NewRows([]string{"id", "name", "isbn", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(1, "Book 1", "1234567890", "2011", "1", "Capybara", "01/01/2002").
		AddRow(2, "Book 2", "0987654321", "2012", "2", "Luffy", "02/02/2002")

	authorRow := sqlMock.NewRows([]string{"author_id", "author_name", "author_birthday"}).
		AddRow("1", "Capybara", "01/01/2002")

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	book := domain.Book{
		Name: "Capybara",
		ISBN: "1234567890",
		Authors: []domain.Author{
			{
				Name: "Capybara Writer",
			},
		},
		PublishYear: 2024,
	}

	fmt.Println("book ", book)
	fmt.Println("test", testBookRepository)

	mock.ExpectQuery("^SELECT \\* FROM author WHERE name = \\$1$").
		WithArgs("Capybara").
		WillReturnRows(authorRow)

	mock.ExpectExec("^WITH new_book AS \\( INSERT INTO book\\(isbn, name, publish_year\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id \\) INSERT INTO book_author\\(book_id, author_id\\) VALUES \\(\\(SELECT id FROM new_book\\), 1\\);$").
		WithArgs(book.ISBN, book.Name, book.PublishYear).
		WillReturnResult(sqlMock.NewResult(1, 1))

	_, err = testBookRepository.CreateBook(book, []string{"Capybara"})
	if err != nil {
		t.Errorf("Error creating book: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
