package router

import (
	"github.com/gorilla/mux"
	"github.com/hoaibao/book-management/handler"
)

func SetMainRouter() *mux.Router {
	return mux.NewRouter()
}

func SetBookRouter(bookHandler *handler.BookHandler, mainRouter *mux.Router) {
	bookRouter := mainRouter.PathPrefix("/books").Subrouter()

	bookRouter.HandleFunc("", bookHandler.GetAllBooksHandler).Methods("GET")
	bookRouter.HandleFunc("/{bookId}", bookHandler.GetBookByIdHandler).Methods("GET")
	bookRouter.HandleFunc("", bookHandler.CreateBookHandler).Methods("POST")
	bookRouter.HandleFunc("/{bookId}", bookHandler.DeleteBookByIdHandler).Methods("DELETE")
	bookRouter.HandleFunc("", bookHandler.DeleteMultipleBookByIdHandler).Methods("DELETE")
	bookRouter.HandleFunc("/{bookId}", bookHandler.UpdateBookByIdHandler).Methods("PUT")
	bookRouter.HandleFunc("", bookHandler.UpdateMultipleBookByIdHandler).Methods("PUT")
}
