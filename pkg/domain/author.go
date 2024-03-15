package domain

type Author struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	BirthDay string `json:"birthDay"`
}
