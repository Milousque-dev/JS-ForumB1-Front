package handlers

import (
	"html/template"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	t, err := template.ParseFiles("./templates/404.tmpl")
	if err != nil {
		http.Error(w, "404 - Page introuvable", http.StatusNotFound)
		return
	}
	t.Execute(w, nil)
}

func InternalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	t, err := template.ParseFiles("./templates/500.tmpl")
	if err != nil {
		http.Error(w, "500 - Erreur interne", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
