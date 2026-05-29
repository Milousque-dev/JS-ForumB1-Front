package models

type TemplateData struct {
	Username string
	Posts []Post
	IsLogged bool
	Error string
}