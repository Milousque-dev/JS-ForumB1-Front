package handlers

import (
	"net/http"
	"fmt"
	"strings"

)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/posts/")
	fmt.Println("Post numéro: ", id)

	if r.Method == http.MethodGet {
		RenderTemplate(w, "post.tmpl", nil)
	}
}