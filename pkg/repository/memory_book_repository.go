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

func (r *MemoryBookRepository) GetAllBooks(isbn, author, fromValue, toValue string) ([]domain.Book, error) {
	r.books = make(map[int]domain.Book)
	result := make([]domain.Book, 0, len(r.books))

	sqlStatement := `
	SELECT 
		b.*, a.* 
	FROM book b 
	JOIN book_author ba ON b.id = ba.book_id
	JOIN author a ON ba.author_id = a.id
	`

	var columnName []string
	var columnValue []string

	if isbn != "" {
		columnName = append(columnName, "b.isbn")
		columnValue = append(columnValue, isbn)
	}

	if author != "" {
		columnName = append(columnName, "a.name")
		columnValue = append(columnValue, author)
	}

	if fromValue != "" && toValue != "" {
		columnName = append(columnName, "b.publish_year")
		columnValue = append(columnValue, fmt.Sprintf("%s %s", fromValue, toValue))
	}

	for i := range columnName {
		if i == 0 {
			sqlStatement += " WHERE "
		}
		if columnName[i] == "b.publish_year" {
			sqlStatement += fmt.Sprintf("b.publish_year >= %s AND b.publish_year <= %s", fromValue, toValue)
		} else {
			sqlStatement += fmt.Sprintf("%s = '%s'", columnName[i], columnValue[i])
		}
		if i < len(columnName)-1 {
			sqlStatement += " AND "
		}
	}

	LogMessage("[SQL]", sqlStatement)
	rows, err := r.DB.Query(sqlStatement)
	CheckError(err, "Error while querying the database")

	for rows.Next() {
		var book domain.Book
		var author domain.Author
		err := rows.Scan(&book.Id, &book.ISBN, &book.Name, &book.PublishYear, &author.Id, &author.Name, &author.BirthDay)
		CheckError(err, "Error while scanning row")
		LogMessage(book)
		if existBook, isExistBook := r.books[book.Id]; isExistBook {
			existBook.Authors = append(existBook.Authors, author)
			r.books[book.Id] = existBook
		} else {
			book.Authors = append(book.Authors, author)
			r.books[book.Id] = book
		}
	}

	for _, value := range r.books {
		result = append(result, value)
	}

	return result, nil
}

func (r *MemoryBookRepository) GetBookById(id int) (domain.Book, error) {
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
		b.*, a.* 
	FROM book b 
	JOIN book_author ba ON b.id = ba.book_id
	JOIN author a ON ba.author_id = a.id
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
		book.Authors = append(book.Authors, author)
	}
	LogMessage(book)
	return book, nil
}

func (r *MemoryBookRepository) CreateBook(book domain.Book, authors []string) (domain.Book, error) {
	var authorSlice []domain.Author
	var unKnowAuthor []string
	for _, value := range authors {
		author, err := memoryAuthorRepository.GetAuthorByName(value)
		CheckError(err, "Not found author")
		if author.Id != -1 {
			authorSlice = append(authorSlice, author)
		} else {
			unKnowAuthor = append(unKnowAuthor, value)
		}
	}
	if len(unKnowAuthor) == 0 {
		insertBookAuthorStatement := "INSERT INTO book_author(book_id, author_id) VALUES "
		bookId := "(SELECT id FROM new_book)"
		for index, value := range authorSlice {
			insertBookAuthorStatement += fmt.Sprintf("(%s, %d)", bookId, value.Id)
			if index < len(authorSlice)-1 {
				insertBookAuthorStatement += ", "
			}
		}
		sqlStatement := fmt.Sprintf(`
		WITH new_book AS (
			INSERT INTO book(isbn, name, publish_year) 
			VALUES ($1, $2, $3)
			RETURNING id
		)

		%s;
		`, insertBookAuthorStatement)
		fmt.Println("SQL: ", sqlStatement)
		LogMessage(sqlStatement)
		result, err := r.DB.Exec(sqlStatement, book.ISBN, book.Name, book.PublishYear)
		CheckError(err, "Can't insert database")
		numberOfRowsAffected, err := result.RowsAffected()
		CheckError(err, "Can't get number of rows affected")
		LogMessage("Number of rows affected:", numberOfRowsAffected)
		LogMessage(book)
		return book, nil
	}

	insertAuthorStatement := "INSERT INTO author(name) VALUES "
	for index, authorName := range unKnowAuthor {
		insertAuthorStatement += fmt.Sprintf("('%s')", authorName)
		if index < len(unKnowAuthor)-1 {
			insertAuthorStatement += ", "
		}
	}
	insertAuthorStatement += "RETURNING id;"
	LogMessage(insertAuthorStatement)
	rows, err := r.DB.Query(insertAuthorStatement)
	CheckError(err, "Can't get id")
	var authorIds []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		CheckError(err, "Error while scanning row")
		LogMessage(book)
		authorIds = append(authorIds, id)
	}
	fmt.Println(authorIds)

	insertBookAuthorStatement := "INSERT INTO book_author(book_id, author_id) VALUES "
	bookId := "(SELECT id FROM new_book)"

	for index, value := range authorSlice {
		insertBookAuthorStatement += fmt.Sprintf("(%s, %d)", bookId, value.Id)
		if index < len(authorSlice)-1 {
			insertBookAuthorStatement += ", "
		}
	}

	for index, value := range authorIds {
		if len(authorSlice) > 0 {
			insertBookAuthorStatement += ", "
		}
		insertBookAuthorStatement += fmt.Sprintf("(%s, %d)", bookId, value)
		if index < len(authorIds)-1 {
			insertBookAuthorStatement += ", "
		}
	}

	sqlStatement := fmt.Sprintf(`
		WITH new_book AS (
			INSERT INTO book(isbn, name, publish_year) 
			VALUES ($1, $2, $3)
			RETURNING id
		)

		%s;
		`, insertBookAuthorStatement)

	LogMessage(sqlStatement)
	result, err := r.DB.Exec(sqlStatement, book.ISBN, book.Name, book.PublishYear)
	CheckError(err, "Can't insert database")
	numberOfRowsAffected, err := result.RowsAffected()
	CheckError(err, "Can't get number of rows affected")
	LogMessage("Number of rows affected:", numberOfRowsAffected)
	LogMessage(book)
	return book, nil
}

func (r *MemoryBookRepository) DeleteBookById(bookId int) (domain.Book, error) {
	book, err := r.GetBookById(bookId)
	CheckError(err, "Book not found")

	sqlStatement := `
	BEGIN TRANSACTION;
	DELETE FROM book_author where book_id = $1;
	DELETE FROM book WHERE id = $2;
	COMMIT;
	`
	LogMessage(sqlStatement)

	tx, err := r.DB.Begin()
	CheckError(err, "Error transaction")

	totalRows := 0

	result, err := tx.Exec("DELETE FROM book_author where book_id = $1", bookId)
	if err != nil {
		tx.Rollback()
	}
	CheckError(err, "Error deleting book_author")
	rows, err := result.RowsAffected()
	totalRows += int(rows)
	CheckError(err, "Error getting row affected")

	result, err = tx.Exec("DELETE FROM book WHERE id = $1", bookId)
	if err != nil {
		tx.Rollback()
	}
	CheckError(err, "Error deleting book")
	rows, err = result.RowsAffected()
	totalRows += int(rows)
	CheckError(err, "Error getting row affected")

	err = tx.Commit()
	CheckError(err, "Error committing transaction")

	LogMessage("Number of rows affected:", totalRows)
	LogMessage(book)
	delete(r.books, bookId)
	return book, nil
}

func (r *MemoryBookRepository) UpdateBookById(bookId int, bookData map[string][]string) (domain.Book, error) {
	existBook, err := r.GetBookById(bookId)
	CheckError(err, "Book not found")
	LogMessage(existBook)

	columnsToUpdate := make([]string, 0, len(bookData))
	newValues := make([]string, 0, len(bookData))

	var authorArr []string
	authorArr = append(authorArr, bookData["author"]...)

	for key, value := range bookData {
		if key == "publishYear" {
			key = "publish_year"
		}

		switch key {
		case "name":
			existBook.Name = value[0]
		case "isbn":
			existBook.ISBN = value[0]
		case "publish_year":
			publishYearInt, err := strconv.Atoi(value[0])
			CheckError(err, "Can't not parse int")
			existBook.PublishYear = publishYearInt
		}

		if key != "author" {
			columnsToUpdate = append(columnsToUpdate, key)
			newValues = append(newValues, value[0])
		}
	}
	fmt.Println("authorArr: ", len(authorArr), authorArr)
	fmt.Println("columnsToUpdate: ", len(columnsToUpdate), columnsToUpdate)
	fmt.Println("newValueshorArr: ", len(newValues), newValues)

	updateBookStatement := "UPDATE book SET "
	for i := 0; i < len(columnsToUpdate); i++ {
		updateBookStatement += fmt.Sprintf("%s = '%s'", columnsToUpdate[i], newValues[i])

		if i < len(columnsToUpdate)-1 {
			updateBookStatement += ", "
		}
	}

	updateBookStatement += " WHERE id = $1"

	if len(authorArr) == 0 {
		LogMessage(updateBookStatement)
		result, err := r.DB.Exec(updateBookStatement, bookId)
		CheckError(err, "Can't update database ")
		LogMessage(result)

		return existBook, nil
	}

	// ---------------------- TEST HERE ------------------------

	var authorObjSlice []domain.Author
	var unKnowAuthor []string
	for _, authorName := range authorArr {
		author, err := memoryAuthorRepository.GetAuthorByName(authorName)
		CheckError(err, "Not found author")
		if author.Id != -1 {
			authorObjSlice = append(authorObjSlice, author)
		} else {
			unKnowAuthor = append(unKnowAuthor, authorName)
		}
	}
	existBook.Authors = authorObjSlice
	fmt.Println("AuthorName: ", authorArr)
	fmt.Println("authorObjSlice: ", authorObjSlice)
	fmt.Println("UnknowAuthor: ", unKnowAuthor)

	totalRows := 0

	var authorIds []int
	if len(unKnowAuthor) > 0 {
		fmt.Println("UnknowAuthor: ", unKnowAuthor, " Unknown Length > 0")
		insertAuthorStatement := "INSERT INTO author(name) VALUES "
		for index, authorName := range unKnowAuthor {
			insertAuthorStatement += fmt.Sprintf("('%s')", authorName)
			if index < len(unKnowAuthor)-1 {
				insertAuthorStatement += ", "
			}
		}
		insertAuthorStatement += "RETURNING id;"
		LogMessage(insertAuthorStatement)
		rows, err := r.DB.Query(insertAuthorStatement)
		CheckError(err, "Can't get id")
		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			CheckError(err, "Error while scanning row")
			authorIds = append(authorIds, id)
		}
	}

	deleteBookAuthorStatement := fmt.Sprintf(`DELETE FROM book_author WHERE book_id = %d`, bookId)
	insertBookAuthorStatement := "INSERT INTO book_author(book_id, author_id) VALUES "

	for index, authorObj := range authorObjSlice {
		insertBookAuthorStatement += fmt.Sprintf("(%d, %d)", bookId, authorObj.Id)
		if index < len(authorObjSlice)-1 {
			insertBookAuthorStatement += ", "
		}
	}

	for index, value := range authorIds {
		if len(authorObjSlice) > 0 {
			insertBookAuthorStatement += ", "
		}
		insertBookAuthorStatement += fmt.Sprintf("(%d, %d)", bookId, value)
		if index < len(authorIds)-1 {
			insertBookAuthorStatement += ", "
		}
	}

	tx, err := r.DB.Begin()
	CheckError(err, "Error beginning transaction")

	LogMessage(deleteBookAuthorStatement)
	result, err := tx.Exec(deleteBookAuthorStatement)
	if err != nil {
		tx.Rollback()
	}
	CheckError(err, "Error deleting existing associations")
	row, err := result.RowsAffected()
	CheckError(err, "Cant not get row affected")
	totalRows += int(row)

	LogMessage(insertBookAuthorStatement)
	result, err = tx.Exec(insertBookAuthorStatement)
	if err != nil {
		tx.Rollback()
	}
	CheckError(err, "Error inserting associations")
	row, err = result.RowsAffected()
	CheckError(err, "Cant not get row affected")
	totalRows += int(row)

	numberOfRowsAffected, err := result.RowsAffected()
	CheckError(err, "Can't get number of rows affected")
	LogMessage("Number of rows affected:", numberOfRowsAffected)

	err = tx.Commit()
	CheckError(err, "Error committing transaction")

	return existBook, nil
}
