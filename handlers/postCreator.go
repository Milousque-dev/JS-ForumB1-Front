package handlers

import (
	"fmt"
	"forum/fake"
	"net/http"
	"strings"
)

func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	_, isLogged := fake.GetCurrentUser(r)

	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	RenderTemplate(w, "postcreate.tmpl", nil)
}

func PostCreator(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))

	if title == "" || content == "" {
		http.Error(w, "Erreur : contenu vide", http.StatusBadRequest)
		return
	}

	fmt.Println(title, content)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

