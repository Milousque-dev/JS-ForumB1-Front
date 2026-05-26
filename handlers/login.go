package handlers

import (
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "login.tmpl", nil)
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println(email, password) //envoyer dans la bdd en vrai :)

	// if (bdd response = ok) {
	http.Redirect(w, r, "/", http.StatusSeeOther) // redirige vers accueil

}