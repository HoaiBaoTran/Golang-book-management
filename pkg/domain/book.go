package domain

type Book struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ISBN        string `json:"isbn"`
	PublishYear int    `json:"publishYear"`
	Author      Author `json:"author"`
}
