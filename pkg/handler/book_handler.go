package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hoaibao/book-management/pkg/domain"
	"github.com/hoaibao/book-management/pkg/service"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) GetAllBooksHandlerVersion2(w http.ResponseWriter, r *http.Request) {
	fromValue := r.URL.Query().Get("from")
	toValue := r.URL.Query().Get("to")

	books, err := h.bookService.GetAllBooksVersion2(fromValue, toValue)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error retrieving books", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) GetBookByIdHandlerVersion2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		http.Error(w, "Invalid book id", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.GetBookByIdVersion2(bookId)
	if err != nil {
		http.Error(w, "Error retrieving book", http.StatusInternalServerError)
	}

	if book == (domain.Book{}) {
		http.Error(w, "Book not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) CreateBookHandlerVersion2(w http.ResponseWriter, r *http.Request) {
	var bookDataList []map[string]string
	var response []domain.Book
	if err := json.NewDecoder(r.Body).Decode(&bookDataList); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, bookData := range bookDataList {
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

		book, err := h.bookService.CreateBookVersion2(name, isbn, author, publishYearInt)
		if err != nil {
			http.Error(w, "Error creating book", http.StatusInternalServerError)
			return
		}
		response = append(response, book)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BookHandler) DeleteMultipleBookByIdHandlerVersion2(w http.ResponseWriter, r *http.Request) {
	var bookDataId map[string][]int
	var response []domain.Book
	if err := json.NewDecoder(r.Body).Decode(&bookDataId); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	bookIdSlice := bookDataId["data"]
	for _, bookId := range bookIdSlice {
		book, err := h.bookService.DeleteBookByIdVersion2(bookId)
		if err != nil {
			http.Error(w, "Error deleting book", http.StatusInternalServerError)
		}
		response = append(response, book)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BookHandler) DeleteBookByIdHandlerVersion2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId, err := strconv.Atoi(vars["bookId"])
	if err != nil {
		http.Error(w, "Invalid book id", http.StatusBadRequest)
		return
	}

	book, err := h.bookService.DeleteBookByIdVersion2(bookId)
	if err != nil {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBookByIdHandlerVersion2(w http.ResponseWriter, r *http.Request) {
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

	book, err := h.bookService.UpdateBookByIdVersion2(bookId, bookData)
	if err != nil {
		log.Fatal("Update fail", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateMultipleBookByIdHandlerVersion2(w http.ResponseWriter, r *http.Request) {
	var bookDataList []map[string]string
	var response []domain.Book
	if err := json.NewDecoder(r.Body).Decode(&bookDataList); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, bookData := range bookDataList {
		if bookId, isContainsBookId := bookData["id"]; isContainsBookId {
			bookIdInt, err := strconv.Atoi(bookId)
			delete(bookData, "id")
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			book, err := h.bookService.UpdateBookByIdVersion2(bookIdInt, bookData)
			if err != nil {
				http.Error(w, "Update failed", http.StatusBadRequest)
				return
			}
			response = append(response, book)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
