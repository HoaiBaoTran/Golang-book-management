package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hoaibao/book-management/service"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	fromValue := r.URL.Query().Get("from")
	toValue := r.URL.Query().Get("to")

	books, err := h.bookService.GetAllBooks(fromValue, toValue)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error retrieving books", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) GetBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		http.Error(w, "Invalid book id", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetBookById(bookId)
	if err != nil {
		http.Error(w, "Error retrieving book", http.StatusInternalServerError)
	}

	if book == nil {
		http.Error(w, "Book not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var bookData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&bookData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	name, nameExists := bookData["name"]
	isbn, isbnExists := bookData["isbn"]
	author, authorExists := bookData["author"]
	publishYear, publishYearExists := bookData["publishYear"]

	if !nameExists || !isbnExists || !authorExists || !publishYearExists {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	publishYearInt, err := strconv.Atoi(publishYear)
	if err != nil {
		http.Error(w, "Invalid publish year", http.StatusBadRequest)
	}

	book, err := h.bookService.CreateBook(name, isbn, author, publishYearInt)
	if err != nil {
		http.Error(w, "Error creating book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) CreateMultipleBookHandler(w http.ResponseWriter, r *http.Request) {
	var books []map[string]string
	json.NewDecoder(r.Body).Decode(&books)
	fmt.Println(books)
}

func (h *BookHandler) DeleteBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		http.Error(w, "Invalid book id", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.DeleteBookById(bookId)
	if err != nil {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBookByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		http.Error(w, "Invalid book id", http.StatusBadRequest)
	}

	var bookData map[string]string
	err = json.NewDecoder(r.Body).Decode(&bookData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.UpdateBookById(bookId, bookData)
	if err != nil {
		log.Fatal("Update fail", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)

}
