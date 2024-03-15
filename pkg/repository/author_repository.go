package repository

import "github.com/hoaibao/book-management/pkg/domain"

type AuthorRepository interface {
	GetAuthorById(id int) (domain.Author, error)
}
