package repository

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/hoaibao/book-management/pkg/database"
	"github.com/hoaibao/book-management/pkg/domain"
	"github.com/hoaibao/book-management/pkg/logger"
	goDotEnv "github.com/joho/godotenv"
)

var (
	MyLogger               = logger.InitLogger()
	memoryAuthorRepository = NewMemoryAuthorRepository()
)

type MemoryBookRepository struct {
	books map[int]domain.Book
	DB    *sql.DB
}

func CheckError(err error, message string) {
	if err != nil {
		MyLogger.ConsoleLogger.Error(message, err)
		MyLogger.FileLogger.Error(message, err)
	}
}

func LogMessage(args ...interface{}) {
	MyLogger.ConsoleLogger.Infoln(args)
	MyLogger.FileLogger.Infoln(args)
}

func NewMemoryBookRepository() *MemoryBookRepository {
	err := goDotEnv.Load(".env")
	CheckError(err, "Can't load value from .env")

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	db, err := database.NewConnection(config)
	CheckError(err, "Can't connect database")

	return &MemoryBookRepository{
		books: make(map[int]domain.Book, 0),
		DB:    db,
	}
}

func (r *MemoryBookRepository) GetAllBooksVersion2(fromValue, toValue string) ([]domain.Book, error) {
	result := make([]domain.Book, 0, len(r.books))

	sqlStatement := `
	SELECT 
		b.id, b.isbn, b.name, b.publish_year, a.id, a.name, a.birth_day 
	FROM book b 
	JOIN author a 
	ON b.author_id = a.id
	`

	var rows *sql.Rows
	var err error

	if fromValue != "" && toValue != "" {
		sqlStatement += " WHERE publish_year >= $1 AND publish_year <= $2"
		LogMessage("[SQL]", sqlStatement)
		rows, err = r.DB.Query(sqlStatement, fromValue, toValue)
	} else {
		LogMessage("[SQL]", sqlStatement)
		rows, err = r.DB.Query(sqlStatement)
	}
	CheckError(err, "Error while querying the database")
	defer rows.Close()

	for rows.Next() {
		var book domain.Book
		var author domain.Author
		err := rows.Scan(&book.Id, &book.ISBN, &book.Name, &book.PublishYear, &author.Id, &author.Name, &author.BirthDay)
		CheckError(err, "Error while scanning row")
		LogMessage(book)
		book.Author = author
		r.books[book.Id] = book
		result = append(result, book)
	}
	return result, nil
}

func (r *MemoryBookRepository) GetBookByIdVersion2(id int) (domain.Book, error) {
	if len(r.books) > 0 {
		book, exist := r.books[id]
		if !exist {
			return domain.Book{}, nil
		}
		LogMessage(book)
		return book, nil
	}

	sqlStatement := `
	SELECT 
		b.id, b.isbn, b.name, b.publish_year, a.id, a.name, a.birth_day 
	FROM book b 
	JOIN author a 
	ON b.author_id = a.id 
	WHERE b.id = $1
	`
	LogMessage(sqlStatement)
	rows, err := r.DB.Query(sqlStatement, id)
	CheckError(err, "Error while querying the database")

	var book domain.Book
	var author domain.Author
	for rows.Next() {
		err := rows.Scan(&book.Id, &book.ISBN, &book.Name, &book.PublishYear, &author.Id, &author.Name, &author.BirthDay)
		CheckError(err, "Error while scanning row")
		LogMessage(book)
		book.Author = author
	}
	return book, nil
}

func (r *MemoryBookRepository) CreateBookVersion2(book domain.Book) (domain.Book, error) {
	author, err := memoryAuthorRepository.GetAuthorByName(book.Author.Name)
	CheckError(err, "Not Found Author")
	if author.Id != -1 {
		sqlStatement := "INSERT INTO book(isbn, name, publish_year, author_id) VALUES ($1, $2, $3, $4)"
		LogMessage(sqlStatement)
		result, err := r.DB.Exec(sqlStatement, book.ISBN, book.Name, book.PublishYear, author.Id)
		CheckError(err, "Can't insert database")
		numberOfRowsAffected, err := result.RowsAffected()
		CheckError(err, "Can't get number of rows affected")
		LogMessage("Number of rows affected:", numberOfRowsAffected)
		LogMessage(book)
		return book, nil
	}

	sqlStatement := `
	WITH new_author AS (
		INSERT INTO author(name) 
		VALUES ($1) RETURNING id
	) 
	INSERT INTO book(isbn, name, publish_year, author_id)
	VALUES
		($2, $3, $4, (SELECT id FROM new_author))
	`
	LogMessage(sqlStatement)
	result, err := r.DB.Exec(sqlStatement, book.Author.Name, book.ISBN, book.Name, book.PublishYear)
	CheckError(err, "Can't insert database")
	numberOfRowsAffected, err := result.RowsAffected()
	CheckError(err, "Can't get number of rows affected")
	LogMessage("Number of rows affected:", numberOfRowsAffected)
	LogMessage(book)
	return book, nil
}

func (r *MemoryBookRepository) DeleteBookByIdVersion2(bookId int) (domain.Book, error) {
	book, err := r.GetBookByIdVersion2(bookId)
	CheckError(err, "Book not found")

	sqlStatement := "DELETE FROM book WHERE id = $1"
	LogMessage(sqlStatement)
	result, err := r.DB.Exec(sqlStatement, bookId)
	CheckError(err, "Can't delete from database")
	numberOfRowsAffected, err := result.RowsAffected()
	CheckError(err, "Can't get number of rows affected")
	LogMessage("Number of rows affected:", numberOfRowsAffected)
	LogMessage(book)
	delete(r.books, bookId)
	return book, nil
}

func (r *MemoryBookRepository) UpdateBookByIdVersion2(bookId int, bookData map[string]string) (domain.Book, error) {
	existBook, err := r.GetBookByIdVersion2(bookId)
	CheckError(err, "Book not found")
	fmt.Println("-------------------", bookId, existBook)

	columnsToUpdate := make([]string, 0, len(bookData))
	newValues := make([]string, 0, len(bookData))

	author, err := memoryAuthorRepository.GetAuthorByName(bookData["author"])
	CheckError(err, "Not Found Author")

	for key, value := range bookData {
		if key == "publishYear" {
			key = "publish_year"
		}

		if key == "author" {
			key = "author_id"
			value = fmt.Sprintf("%d", author.Id)
		}

		switch key {
		case "name":
			existBook.Name = value
		case "isbn":
			existBook.ISBN = value
		case "publish_year":
			publishYearInt, err := strconv.Atoi(value)
			CheckError(err, "Can't not parse int")
			existBook.PublishYear = publishYearInt
		case "author_id":
			existBook.Author = author
		}

		columnsToUpdate = append(columnsToUpdate, key)
		newValues = append(newValues, value)
	}

	if author.Id != -1 {
		sqlStatement := "UPDATE book SET "
		for i := 0; i < len(columnsToUpdate); i++ {
			if columnsToUpdate[i] == "author_id" {
				sqlStatement += fmt.Sprintf("%s = %s", columnsToUpdate[i], newValues[i])
			} else {
				sqlStatement += fmt.Sprintf("%s = '%s'", columnsToUpdate[i], newValues[i])
			}

			if i < len(columnsToUpdate)-1 {
				sqlStatement += ", "
			}
		}

		sqlStatement += " WHERE id = $1"
		LogMessage(sqlStatement)

		result, err := r.DB.Exec(sqlStatement, bookId)
		CheckError(err, "Can't update database ")
		LogMessage(result)

		return existBook, nil
	}

	sqlStatement := `
	WITH new_author AS (
		INSERT INTO author(name) 
		VALUES ($1) RETURNING id::int
	) 
	UPDATE book b SET `
	for i := 0; i < len(columnsToUpdate); i++ {
		if columnsToUpdate[i] == "author_id" {
			sqlStatement += fmt.Sprintf("%s = (SELECT id FROM new_author)", columnsToUpdate[i])
		} else {
			sqlStatement += fmt.Sprintf("%s = '%s'", columnsToUpdate[i], newValues[i])
		}
		if i < len(columnsToUpdate)-1 {
			sqlStatement += ", "
		}
	}

	sqlStatement += " WHERE b.id = $2"
	LogMessage(sqlStatement)

	result, err := r.DB.Exec(sqlStatement, bookData["author"], bookId)
	CheckError(err, "Can't update database ")
	LogMessage(result)

	existBook.Author = domain.Author{
		Name: bookData["author"],
	}

	return existBook, nil
}
