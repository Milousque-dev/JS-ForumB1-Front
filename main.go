package main

import (
	"fmt"
	"forum/database"
	"net/http"
	"forum/handlers"
)

const port = ":8080"

func main() {
	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Erreur database: ", err)
		return
	}

	defer db.Close()
	fmt.Println("Database créée et fonctionnelle")

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/posts/", handlers.PostHandler)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Serveur lancé sur (http://localhost" + port + ")")
	
	err = http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Println("erreur serveur:", err)
	}
}