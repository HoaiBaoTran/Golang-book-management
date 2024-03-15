package repository

import (
	"database/sql"
	"os"

	"github.com/hoaibao/book-management/pkg/database"
	"github.com/hoaibao/book-management/pkg/domain"
	goDotEnv "github.com/joho/godotenv"
)

type MemoryAuthorRepository struct {
	authors map[int]domain.Author
	DB      *sql.DB
}

func NewMemoryAuthorRepository() *MemoryAuthorRepository {
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

	return &MemoryAuthorRepository{
		authors: make(map[int]domain.Author, 0),
		DB:      db,
	}
}

func (authorRepository *MemoryAuthorRepository) GetAuthorById(id int) (domain.Author, error) {
	if len(authorRepository.authors) != 0 {
		author, exist := authorRepository.authors[id]
		if !exist {
			return domain.Author{}, nil
		}
		LogMessage(author)
		return author, nil
	}

	sqlStatement := "SELECT * FROM author WHERE id = $1"
	LogMessage(sqlStatement)
	rows, err := authorRepository.DB.Query(sqlStatement, id)
	CheckError(err, "Error while querying the database")

	var author domain.Author
	for rows.Next() {
		err := rows.Scan(&author.Id, &author.Name, &author.BirthDay)
		CheckError(err, "Error while scanning row")
		LogMessage(author)
	}
	return author, nil
}

func (authorRepository *MemoryAuthorRepository) GetAuthorByName(name string) (domain.Author, error) {
	if len(authorRepository.authors) != 0 {
		for _, author := range authorRepository.authors {
			if author.Name == name {
				LogMessage(author)
				return author, nil
			}
		}
	}
	sqlStatement := "SELECT * FROM author WHERE name = $1"
	LogMessage(sqlStatement)
	rows, err := authorRepository.DB.Query(sqlStatement, name)
	CheckError(err, "Error while querying the database")

	var author domain.Author
	if rows.Next() {
		err := rows.Scan(&author.Id, &author.Name, &author.BirthDay)
		CheckError(err, "Error while scanning row")
		LogMessage(author)
		return author, nil
	}
	return domain.Author{
		Id: -1,
	}, nil
}
