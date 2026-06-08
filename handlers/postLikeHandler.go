package handlers

import (
	"forum/database"
	"forum/fake"
	"net/http"
	"strconv"
)

func PostLikeHandler(w http.ResponseWriter, r *http.Request) {
	togglePostReaction(w, r, 1)
}

func PostDislikeHandler(w http.ResponseWriter, r *http.Request) {
	togglePostReaction(w, r, -1)
}

func togglePostReaction(w http.ResponseWriter, r *http.Request, value int) {
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

	database.TogglePostLike(user.ID, postID, value)

	http.Redirect(w, r, "/posts/"+postIDStr, http.StatusSeeOther)
}
