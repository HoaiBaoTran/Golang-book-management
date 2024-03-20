package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	sqlMock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/hoaibao/book-management/pkg/domain"
	"github.com/hoaibao/book-management/pkg/handler"
	"github.com/hoaibao/book-management/pkg/repository"
	"github.com/hoaibao/book-management/pkg/service"
)

func TestBookHandler_GetAllBookHandler(t *testing.T) {

	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlMock.NewRows([]string{"id", "isbn", "name", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(1, "1234567890", "Book 1", "2011", 1, "Capybara", "01/01/2002").
		AddRow(2, "0987654321", "Book 2", "2012", 2, "Luffy", "02/02/2002")

	mock.ExpectQuery("^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id$").
		WillReturnRows(rows)

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	req, err := http.NewRequest("GET", "/api/v3/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a mock HTTP handler function for your API route
	handler := http.HandlerFunc(bookHandler.GetAllBooksHandler)

	// Serve the HTTP request to the mock handler
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":1,"isbn":"1234567890","name":"Book 1","authors":[{"id":1,"name":"Capybara","birthDay":"01/01/2002"}],"publishYear":2011},{"id":2,"isbn":"0987654321","name":"Book 2","authors":[{"id":2,"name":"Luffy","birthDay":"02/02/2002"}],"publishYear":2012}]`
	actual := rr.Body.String()

	var expectedBooks []domain.Book
	var actualBooks []domain.Book
	if err := json.Unmarshal([]byte(expected), &expectedBooks); err != nil {
		t.Errorf("error unMarshalling expected JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(actual), &actualBooks); err != nil {
		t.Errorf("error unMarshalling actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedBooks, actualBooks) {
		t.Errorf("expected books and actual books are not equal got %v want %v", actualBooks, expectedBooks)
	}
}

func TestBookHandler_GetBookHandlerWithFilter(t *testing.T) {

	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	rows := sqlMock.NewRows([]string{"id", "isbn", "name", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(1, "1234567890", "Book 1", "2011", 1, "Capybara", "01/01/2002").
		AddRow(2, "0987654321", "Book 2", "2012", 2, "Luffy", "02/02/2002")

	mock.ExpectQuery("^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id WHERE b.publish_year >= 2010 AND b.publish_year <= 2013$").
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/api/v3/books?from=2010&to=2013", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a mock HTTP handler function for your API route
	handler := http.HandlerFunc(bookHandler.GetAllBooksHandler)

	// Serve the HTTP request to the mock handler
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":1,"isbn":"1234567890","name":"Book 1","authors":[{"id":1,"name":"Capybara","birthDay":"01/01/2002"}],"publishYear":2011},{"id":2,"isbn":"0987654321","name":"Book 2","authors":[{"id":2,"name":"Luffy","birthDay":"02/02/2002"}],"publishYear":2012}]`
	actual := rr.Body.String()

	var expectedBooks []domain.Book
	var actualBooks []domain.Book
	if err := json.Unmarshal([]byte(expected), &expectedBooks); err != nil {
		t.Errorf("error unmarshalling expected JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(actual), &actualBooks); err != nil {
		t.Errorf("error unmarshalling actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedBooks, actualBooks) {
		t.Errorf("expected books and actual books are not equal got %v want %v", actualBooks, expectedBooks)
	}
}

func TestBookHandler_GetBookByIdHandler(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlMock.NewRows([]string{"id", "isbn", "name", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(2, "0987654321", "Book 2", "2012", 2, "Luffy", "02/02/2002")

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	req, err := http.NewRequest("GET", "/api/v3/books/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery("^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id WHERE b.id = \\$1$").
		WithArgs(2).
		WillReturnRows(rows)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v3/books/{bookId}", bookHandler.GetBookByIdHandler)

	// Serve the HTTP request to the mock handler
	r.ServeHTTP(rr, req)
	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":2,"isbn":"0987654321","name":"Book 2","authors":[{"id":2,"name":"Luffy","birthDay":"02/02/2002"}],"publishYear":2012}`
	actual := rr.Body.String()

	var expectedBooks domain.Book
	var actualBooks domain.Book
	if err := json.Unmarshal([]byte(expected), &expectedBooks); err != nil {
		t.Errorf("error unmarshalling expected JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(actual), &actualBooks); err != nil {
		t.Errorf("error unmarshalling actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedBooks, actualBooks) {
		t.Errorf("expected books and actual books are not equal got %v want %v", actualBooks, expectedBooks)
	}
}

func TestBookHandler_CreateBookHandler(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	rows := sqlMock.NewRows([]string{"author_id", "author_name", "author_birthday"}).
		AddRow(1, "Author 1", "01/01/2002")

	book := `[{"name":["Book 1"],"isbn":["1234567890"],"author":["Author 1"],"publishYear":["2010"]}]`
	fmt.Println("payload: ", strings.NewReader(book))

	mock.ExpectQuery("^SELECT \\* FROM author WHERE name = \\$1$").
		WithArgs("Author 1").
		WillReturnRows(rows)

	mock.ExpectExec("^WITH new_book AS \\( INSERT INTO book\\(isbn, name, publish_year\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id \\) INSERT INTO book_author\\(book_id, author_id\\) VALUES \\(\\(SELECT id FROM new_book\\), 1\\);$").
		WithArgs("1234567890", "Book 1", 2010).
		WillReturnResult(sqlMock.NewResult(1, 1))

	req, err := http.NewRequest("POST", "/api/v3/books", strings.NewReader(book))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(bookHandler.CreateBookHandler)

	handler.ServeHTTP(rr, req)
	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	fmt.Println("Body: ", rr.Body)

	var returnedBook []domain.Book
	err = json.NewDecoder(rr.Body).Decode(&returnedBook)
	if err != nil {
		fmt.Println("body: ", rr.Body)
		t.Errorf("error parsing response body: %v", err)
	}

	expected := `[{"id":0,"name":"Book 1","isbn":"1234567890","authors":[{"id":0,"name":"Author 1","birthDay":""}],"publishYear":2010}]`
	var expectedBooks []domain.Book
	if err := json.Unmarshal([]byte(expected), &expectedBooks); err != nil {
		t.Errorf("error unmarshalling expected JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedBooks, returnedBook) {
		t.Errorf("expected books and actual books are not equal got %v want %v", returnedBook, expectedBooks)
	}
}

func TestBookHandler_DeleteBookByIdHandler(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlMock.NewRows([]string{"id", "isbn", "name", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(2, "0987654321", "Book 2", "2012", 2, "Luffy", "02/02/2002")

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	req, err := http.NewRequest("DELETE", "/api/v3/books/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery("^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id WHERE b.id = \\$1$").
		WithArgs(2).
		WillReturnRows(rows)

	mock.ExpectExec("^DELETE FROM book WHERE id = \\$1$").
		WithArgs(2).
		WillReturnResult(sqlMock.NewResult(0, 1))

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v3/books/{bookId}", bookHandler.GetBookByIdHandler)

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":2,"name":"Book 2","isbn":"0987654321","authors":[{"id":2, "name":"Luffy", "birthDay":"02/02/2002"}],"publishYear":2012}`
	actual := rr.Body.String()

	var expectedBooks domain.Book
	var actualBooks domain.Book
	if err := json.Unmarshal([]byte(expected), &expectedBooks); err != nil {
		t.Errorf("error unmarshalling expected JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(actual), &actualBooks); err != nil {
		t.Errorf("error unmarshalling actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedBooks, actualBooks) {
		t.Errorf("expected books and actual books are not equal got %v want %v", actualBooks, expectedBooks)
	}
}

func TestBookHandler_UpdateBookByIdHandler(t *testing.T) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlMock.NewRows([]string{"id", "isbn", "name", "publish_year", "author_id", "author_name", "author_birthday"}).
		AddRow(2, "0987654321", "Book 2", "2012", 2, "Luffy", "02/02/2002")

	authorRow := sqlMock.NewRows([]string{"author_id", "author_name", "author_birthday"}).
		AddRow(1, "Capybara", "01/01/2002")

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	book := `{"name":["Luffy"],"author":["Capybara"]}`

	fmt.Println("payload: ", strings.NewReader(book))
	req, err := http.NewRequest("PUT", "/api/v3/books/2", strings.NewReader(book))
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery("^SELECT b.*, a.* FROM book b JOIN book_author ba ON b.id = ba.book_id JOIN author a ON ba.author_id = a.id WHERE b.id = \\$1$").
		WithArgs(2).
		WillReturnRows(rows)

	mock.ExpectQuery("^SELECT \\* FROM author WHERE name = \\$1$").
		WithArgs("Capybara").
		WillReturnRows(authorRow)

	// mock.ExpectExec("^UPDATE book SET name = 'Luffy', author_id = '1' WHERE id = \\$1$").
	// 	WithArgs(2).
	// 	WillReturnResult(sqlMock.NewResult(1, 1))

	mock.ExpectBegin()

	mock.ExpectExec("^DELETE FROM book_author WHERE book_id = 2$").
		WillReturnResult(sqlMock.NewResult(1, 1))

	mock.ExpectExec("^INSERT INTO book_author\\(book_id, author_id\\) VALUES \\(2, 1\\)$").
		WillReturnResult(sqlMock.NewResult(1, 1))

	mock.ExpectCommit()

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v3/books/{bookId}", bookHandler.UpdateBookByIdHandler)

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":2,"name":"Luffy","isbn":"0987654321","authors":[{"id":1,"name":"Capybara","birthDay":"01/01/2002"}],"publishYear":2012}`
	actual := rr.Body.String()

	var expectedBooks domain.Book
	var actualBooks domain.Book
	if err := json.Unmarshal([]byte(expected), &expectedBooks); err != nil {
		t.Errorf("error unmarshalling expected JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(actual), &actualBooks); err != nil {
		fmt.Println("Actual Book: ", actualBooks)
		t.Errorf("error unmarshalling actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedBooks, actualBooks) {
		t.Errorf("expected books and actual books are not equal got %v want %v", actualBooks, expectedBooks)
	}
}