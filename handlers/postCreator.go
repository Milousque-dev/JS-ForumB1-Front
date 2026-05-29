package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
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

