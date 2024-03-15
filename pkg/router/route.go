package router

import (
	"github.com/gorilla/mux"
	"github.com/hoaibao/book-management/pkg/handler"
)

func SetMainRouter() *mux.Router {
	return mux.NewRouter()
}

func SetBookRouter(bookHandler *handler.BookHandler, mainRouter *mux.Router) {
	bookRouter := mainRouter.PathPrefix("/api/v2/books").Subrouter()

	bookRouter.HandleFunc("", bookHandler.GetAllBooksHandlerVersion2).Methods("GET")
	bookRouter.HandleFunc("/{bookId}", bookHandler.GetBookByIdHandlerVersion2).Methods("GET")
	bookRouter.HandleFunc("", bookHandler.CreateBookHandlerVersion2).Methods("POST")
	bookRouter.HandleFunc("/{bookId}", bookHandler.DeleteBookByIdHandlerVersion2).Methods("DELETE")
	bookRouter.HandleFunc("", bookHandler.DeleteMultipleBookByIdHandlerVersion2).Methods("DELETE")
	bookRouter.HandleFunc("/{bookId}", bookHandler.UpdateBookByIdHandlerVersion2).Methods("PUT")
	bookRouter.HandleFunc("", bookHandler.UpdateMultipleBookByIdHandlerVersion2).Methods("PUT")
}
