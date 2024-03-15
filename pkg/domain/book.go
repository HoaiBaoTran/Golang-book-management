package domain

type Book struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ISBN        string `json:"isbn"`
	Author      string `json:"author"`
	PublishYear int    `json:"publishYear"`
}
