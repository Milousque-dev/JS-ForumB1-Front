package handlers

import (
	"net/http"
	"html/template"
	"forum/fake"
	"forum/models"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles("./templates/" + tmpl)
	if err != nil {
		InternalErrorHandler(w, nil)
		return
	}
	if err = t.Execute(w, data); err != nil {
		InternalErrorHandler(w, nil)
	}
}

func Home( w http.ResponseWriter, r *http.Request) {
	posts := fake.GetAllPosts()

	username, isLogged := fake.GetCurrentUser(r)

	data := models.TemplateData {
		Username: username,
		Posts: posts,
		IsLogged: isLogged,
		Error: "",
	}
	RenderTemplate(w, "index.tmpl", data)
}