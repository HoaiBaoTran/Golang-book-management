package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*sql.DB, error) {
	dataSource := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		log.Fatal("Can't connect database", err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Can't connect database", err)
	}

	fmt.Println("Connect database successfully")
	return db, nil
}
