package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hoaibao/book-management/pkg/database"
	"github.com/hoaibao/book-management/pkg/domain"
	"github.com/hoaibao/book-management/pkg/logger"
	goDotEnv "github.com/joho/godotenv"
)

var myLogger = logger.InitLogger()

type MemoryBookRepository struct {
	books map[int]domain.Book
	DB    *sql.DB
}

func checkError(err error, message string) {
	if err != nil {
		myLogger.ConsoleLogger.Error(message, err)
		myLogger.FileLogger.Error(message, err)
	}
}

func logMessage(args ...interface{}) {
	myLogger.ConsoleLogger.Infoln(args)
	myLogger.FileLogger.Infoln(args)
}

func NewMemoryBookRepository() *MemoryBookRepository {
	err := goDotEnv.Load(".env")
	checkError(err, "Can't load value from .env")

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	db, err := database.NewConnection(config)
	checkError(err, "Can't connect database")

	return &MemoryBookRepository{
		books: make(map[int]domain.Book, 0),
		DB:    db,
	}
}

func (r *MemoryBookRepository) GetAllBooks(fromValue, toValue string) ([]domain.Book, error) {
	result := make([]domain.Book, 0, len(r.books))

	sqlStatement := "SELECT * FROM book"

	var rows *sql.Rows
	var err error

	if fromValue != "" && toValue != "" {
		sqlStatement += " WHERE publish_year >= $1 AND publish_year <= $2"
		logMessage("[SQL]", sqlStatement)
		rows, err = r.DB.Query(sqlStatement, fromValue, toValue)
	} else {
		logMessage("[SQL]", sqlStatement)
		rows, err = r.DB.Query(sqlStatement)
	}
	checkError(err, "Error while querying the database")
	defer rows.Close()

	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.Id, &book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		checkError(err, "Error while scanning row")
		logMessage(book)
		r.books[book.Id] = book
		result = append(result, book)
	}
	return result, nil
}

func (r *MemoryBookRepository) GetBookById(id int) (domain.Book, error) {
	if len(r.books) != 0 {
		book, exist := r.books[id]
		if !exist {
			return domain.Book{}, nil
		}
		logMessage(book)
		return book, nil
	}

	sqlStatement := "SELECT * FROM book WHERE id = $1"
	logMessage(sqlStatement)
	rows, err := r.DB.Query(sqlStatement, id)
	checkError(err, "Error while querying the database")

	var book domain.Book
	for rows.Next() {
		err := rows.Scan(&book.Id, &book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		checkError(err, "Error while scanning row")
		logMessage(book)
	}
	return book, nil
}

func (r *MemoryBookRepository) CreateBook(book domain.Book) (domain.Book, error) {
	sqlStatement := "INSERT INTO book(isbn, name, author, publish_year) VALUES ($1, $2, $3, $4)"
	logMessage(sqlStatement)
	result, err := r.DB.Exec(sqlStatement, book.ISBN, book.Name, book.Author, book.PublishYear)
	checkError(err, "Can't insert database")
	numberOfRowsAffected, err := result.RowsAffected()
	checkError(err, "Can't get number of rows affected")
	logMessage("Number of rows affected:", numberOfRowsAffected)
	logMessage(book)
	return book, nil
}

func (r *MemoryBookRepository) DeleteBookById(bookId int) (domain.Book, error) {
	book, err := r.GetBookById(bookId)
	checkError(err, "Book not found")

	sqlStatement := "DELETE FROM book WHERE id = $1"
	logMessage(sqlStatement)
	result, err := r.DB.Exec(sqlStatement, bookId)
	checkError(err, "Can't delete from database")
	numberOfRowsAffected, err := result.RowsAffected()
	checkError(err, "Can't get number of rows affected")
	logMessage("Number of rows affected:", numberOfRowsAffected)
	logMessage(book)
	delete(r.books, bookId)
	return book, nil
}

func (r *MemoryBookRepository) UpdateBookById(bookId int, bookData map[string]string) (domain.Book, error) {
	existBook, err := r.GetBookById(bookId)
	checkError(err, "Book not found")

	columnsToUpdate := make([]string, 0, len(bookData))
	newValues := make([]string, 0, len(bookData))

	for key, value := range bookData {
		if key == "publishYear" {
			key = "publish_year"
		}
		columnsToUpdate = append(columnsToUpdate, key)
		newValues = append(newValues, value)
		switch key {
		case "name":
			existBook.Name = value
		case "isbn":
			existBook.ISBN = value
		case "author":
			existBook.Author = value
		case "publishYear":
			publishYearInt, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal("Cant update", err)
			}
			existBook.PublishYear = publishYearInt
		default:
			continue
		}
	}

	sqlStatement := "UPDATE book SET "
	for i := 0; i < len(columnsToUpdate); i++ {
		sqlStatement += fmt.Sprintf("%s = '%s'", columnsToUpdate[i], newValues[i])
		if i < len(columnsToUpdate)-1 {
			sqlStatement += ", "
		}
	}

	sqlStatement += " WHERE id = $1"
	logMessage(sqlStatement)

	result, err := r.DB.Exec(sqlStatement, bookId)
	checkError(err, "Can't update database ")
	logMessage(result)

	return existBook, nil
}
