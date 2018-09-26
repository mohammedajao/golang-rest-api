package models

type Post struct {
	id          int
	user_id     int
	likes       int
	title       string
	description string
	slug        string
	category    string
}
