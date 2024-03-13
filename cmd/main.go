package main

import (
	"fmt"
	"net/http"

	"github.com/hoaibao/book-management/handler"
	"github.com/hoaibao/book-management/repository"
	"github.com/hoaibao/book-management/router"
	"github.com/hoaibao/book-management/service"
)

func main() {

	bookRepository := repository.NewMemoryBookRepository()
	bookService := service.NewBookService(bookRepository)
	bookHandler := handler.NewBookHandler(bookService)

	mainRouter := router.SetMainRouter()
	router.SetBookRouter(bookHandler, mainRouter)

	fmt.Println("Starting server at port 8080")
	http.ListenAndServe(":8080", mainRouter)
}
