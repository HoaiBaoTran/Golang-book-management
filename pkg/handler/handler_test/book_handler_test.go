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

	rows := sqlMock.NewRows([]string{"id", "name", "isbn", "author", "publish_year"}).
		AddRow(1, "1234567890", "Book 1", "Author 1", "2011").
		AddRow(2, "0987654321", "Book 2", "Author 2", "2012")

	mock.ExpectQuery("^SELECT \\* FROM book$").
		WillReturnRows(rows)

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	req, err := http.NewRequest("GET", "/api/v1/books", nil)
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

	expected := `[{"id":1,"name":"Book 1","isbn":"1234567890","author":"Author 1","publishYear":2011},{"id":2,"name":"Book 2","isbn":"0987654321","author":"Author 2","publishYear":2012}]`
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
		t.Errorf("expected books and actual books are not equal")
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

	rows := sqlMock.NewRows([]string{"id", "isbn", "name", "author", "publish_year"}).
		AddRow(1, "1234567890", "Book 1", "Author 1", "2011").
		AddRow(2, "0987654321", "Book 2", "Author 2", "2012")

	mock.ExpectQuery("^SELECT \\* FROM book WHERE publish_year >= 2010 AND publish_year <= 2013$").
		WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/api/v1/books?from=2010&to=2013", nil)
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

	expected := `[{"id":1,"name":"Book 1","isbn":"1234567890","author":"Author 1","publishYear":2011},{"id":2,"name":"Book 2","isbn":"0987654321","author":"Author 2","publishYear":2012}]`
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

	rows := sqlMock.NewRows([]string{"id", "isbn", "name", "author", "publish_year"}).
		AddRow(1, "1234567890", "Book 1", "Author 1", "2011").
		AddRow(2, "0987654321", "Book 2", "Author 2", "2012")

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	req, err := http.NewRequest("GET", "/api/v1/books/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery("^SELECT \\* FROM book WHERE id = \\$1$").
		WithArgs(2).
		WillReturnRows(rows)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/books/{bookId}", bookHandler.GetBookByIdHandler)

	// Serve the HTTP request to the mock handler
	r.ServeHTTP(rr, req)
	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":2,"name":"Book 2","isbn":"0987654321","author":"Author 2","publishYear":2012}`
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
		t.Errorf("expected books and actual books are not equal")
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

	mock.ExpectExec("^INSERT INTO book\\(isbn, name, author, publish_year\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)$").
		WithArgs("1234567890", "Book 1", "Author 1", 2010).
		WillReturnResult(sqlMock.NewResult(1, 1))

	book := `[{"id":"1","name":"Book 1","isbn":"1234567890","author":"Author 1","publishYear":"2010"}]`

	fmt.Println("payload: ", strings.NewReader(book))
	req, err := http.NewRequest("POST", "/api/v1/books", strings.NewReader(book))
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

	expected := `[{"id":0,"name":"Book 1","isbn":"1234567890","author":"Author 1","publishYear":2010}]`
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

	rows := sqlMock.NewRows([]string{"id", "isbn", "name", "author", "publish_year"}).
		AddRow(1, "1234567890", "Book 1", "Author 1", "2011").
		AddRow(2, "0987654321", "Book 2", "Author 2", "2012")

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	req, err := http.NewRequest("DELETE", "/api/v1/books/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery("^SELECT \\* FROM book WHERE id = \\$1$").
		WithArgs(2).
		WillReturnRows(rows)

	mock.ExpectExec("^DELETE FROM book WHERE id = \\$1$").
		WithArgs(2).
		WillReturnResult(sqlMock.NewResult(0, 1))

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/books/{bookId}", bookHandler.GetBookByIdHandler)

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":2,"name":"Book 2","isbn":"0987654321","author":"Author 2","publishYear":2012}`
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

	rows := sqlMock.NewRows([]string{"id", "isbn", "name", "author", "publish_year"}).
		AddRow(1, "1234567890", "Book 1", "Author 1", "2011").
		AddRow(2, "0987654321", "Book 2", "Author 2", "2012")

	testBookRepository := repository.NewTestMemoryBookRepository(db)
	bookService := service.NewBookService(testBookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	book := `{"name":"Luffy","author":"Capybara"}`

	fmt.Println("payload: ", strings.NewReader(book))
	req, err := http.NewRequest("PUT", "/api/v1/books/2", strings.NewReader(book))
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery("^SELECT \\* FROM book WHERE id = \\$1$").
		WithArgs(2).
		WillReturnRows(rows)

	mock.ExpectExec("^UPDATE book SET name = 'Luffy', author = 'Capybara' WHERE id = \\$1$").
		WithArgs(2).
		WillReturnResult(sqlMock.NewResult(1, 1))

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/books/{bookId}", bookHandler.UpdateBookByIdHandler)

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"id":2,"name":"Luffy","isbn":"0987654321","author":"Capybara","publishYear":2012}`
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
