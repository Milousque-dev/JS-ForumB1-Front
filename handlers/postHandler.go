package handlers

import (
	"forum/database"
	"forum/fake"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post, ok := fake.GetPostById(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	username, isLogged := fake.GetCurrentUser(r)
	comments := fake.GetCommentByPostID(id)

	datas := models.TemplateData{
		Username: username,
		IsLogged: isLogged,
		Post:     post,
		Comments: comments,
	}

	RenderTemplate(w, "post.tmpl", datas)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	user, isLogged := fake.GetCurrentUserFull(r)
	if !isLogged {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := r.PathValue("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, found := fake.GetPostById(postID)
	if !found {
		http.NotFound(w, r)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Redirect(w, r, "/posts/"+postIDStr, http.StatusSeeOther)
		return
	}

	if err := database.CreateComment(content, user.ID, postID); err != nil {
		http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts/"+postIDStr, http.StatusSeeOther)
}
