package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hoaibao/book-management/database"
	"github.com/hoaibao/book-management/domain"
	"github.com/joho/godotenv"
)

type MemoryBookRepository struct {
	books map[int]domain.Book
	DB    *sql.DB
}

func NewMemoryBookRepository() *MemoryBookRepository {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatal("Can't connect database", err)
	}

	return &MemoryBookRepository{
		books: make(map[int]domain.Book, 0),
		DB:    db,
	}
}

func (r *MemoryBookRepository) GetAllBooks() ([]domain.Book, error) {
	result := make([]domain.Book, 0, len(r.books))

	sqlStatement := "SELECT * FROM book"
	rows, err := r.DB.Query(sqlStatement)
	if err != nil {
		log.Fatal("Error querying the database", err)
	}
	defer rows.Close()

	for rows.Next() {
		var book domain.Book
		err := rows.Scan(&book.Id, &book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		if err != nil {
			log.Fatal("Error scanning row")
		}
		r.books[book.Id] = book
		result = append(result, book)
	}

	return result, nil
}

func (r *MemoryBookRepository) GetBookById(id int) (*domain.Book, error) {
	r.GetAllBooks()
	book, exist := r.books[id]
	if !exist {
		return nil, nil
	}
	return &book, nil
}

func (r *MemoryBookRepository) CreateBook(book *domain.Book) (*domain.Book, error) {
	sqlStatement := "INSERT INTO book(isbn, name, author, publish_year) VALUES ($1, $2, $3, $4)"
	_, err := r.DB.Exec(sqlStatement, book.ISBN, book.Name, book.Author, book.PublishYear)
	if err != nil {
		log.Fatal("Can't insert database", err)
	}
	fmt.Println("INSERT SUCCESSFULLY")
	return book, nil
}

func (r *MemoryBookRepository) DeleteBookById(bookId int) (*domain.Book, error) {

	book, err := r.GetBookById(bookId)
	if err != nil {
		log.Fatal("Book not found", err)
	}

	sqlStatement := "DELETE FROM book WHERE id = $1"
	_, err = r.DB.Exec(sqlStatement, bookId)
	if err != nil {
		log.Fatal("Can't delete from database", err)
	}

	delete(r.books, bookId)
	fmt.Println("DELETE SUCCESSFULLY")
	return book, nil
}

func (r *MemoryBookRepository) UpdateBookById(bookId int, bookData map[string]string) (*domain.Book, error) {
	existBook, err := r.GetBookById(bookId)
	if err != nil {
		log.Fatal("Book not found", err)
	}

	columnsToUpdate := make([]string, 0, len(bookData))
	newValues := make([]string, 0, len(bookData))

	for key, value := range bookData {
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

	_, err = r.DB.Exec(sqlStatement, bookId)
	if err != nil {
		log.Fatal("Can't insert database ", err)
	}

	fmt.Println("UPDATE SUCCESSFULLY")

	return existBook, nil
}
